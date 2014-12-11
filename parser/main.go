package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/haklop/bazooka/commons/matrix"

	log "github.com/Sirupsen/logrus"
	lib "github.com/haklop/bazooka/commons"
)

const (
	SourceFolder      = "/bazooka"
	OutputFolder      = "/bazooka-output"
	MetaFolder        = "/meta"
	BazookaConfigFile = ".bazooka.yml"
	TravisConfigFile  = ".travis.yml"
	MX_ENV_PREFIX     = "env::"
)

func main() {
	log.Info("Starting Parsing Phase")
	// Find either .travis.yml or .bazooka.yml file in the project
	configFile, err := lib.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	// parse the configuration
	config := &lib.Config{}
	err = lib.Parse(configFile, config)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"config": config,
	}).Debug("Configuration parsed")

	// resolve the docker image corresponding to this particular language parser
	image, err := resolveLanguageParser(config.Language)
	if err != nil {
		log.Fatal(err)
	}

	// run the parser image
	langParser := &LanguageParser{
		Image: image,
	}
	err = langParser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// if all went well, the parser should have generated one or more "sub" .bazooka.*.yml files
	// one for each compiler version for example
	//
	// they are also supposed to enrich it with a from attribute corresponding to a base docker image
	// to be used to run the build

	log.Info("Starting Matrix generation")
	files, err := lib.ListFilesWithPrefix(OutputFolder, ".bazooka")
	if err != nil {
		log.Fatal(err)
	}

	// for each of those files (the "sub" .bazooka.*.yml)
	for _, file := range files {
		// parse the damned thing
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			log.Fatal(fmt.Errorf("Error while parsing config file %s: %v", file, err))
		}

		// create a matrix from the environment variables
		// a matrix has N variables (dimensions) where each variable has M values
		// config.Env is a list of key=value strings
		// explodeProps transforms that into a map[string][]string
		// for example ["A=1", "A=2", B="3"]
		// when exploded gets transformed into
		// {A: [1, 2], B: [3]}
		// this matches the matrix layout, so we could store them directly into the matrix
		// but since we would like to be able to extract them later, and to avoid mixing them with a language specific variables (like jdk, go, etc.)
		// explode prefixes the env variables names with a prefix defined in the constant MX_ENV_PREFIX
		// Hence, our matrix is more like: {"env::A": [1, 2], "env::B": [3]}
		mx := matrix.Matrix(explodeProps(config.Env, MX_ENV_PREFIX))

		// extract the "*" part from the .bazooka.*.yml file
		rootCounter := parseCounter(file)
		// for every .bazooka.*.yml file, the language parser is also supposed to have generated a meta/* file
		// which is a simple yml file containing the language specific  matrix variables
		// for example, if the original .bazooka.yml file defined 2 go versions:
		//
		// go:
		// - 1.2.2
		// - 1.3.1
		//
		// the language parser should generate 2 meta files, one for each go version in this format:
		//
		// go: 1.2.2
		//
		// and
		//
		// go: 1.3.1
		rootMetaFile := fmt.Sprintf("%s/%s", MetaFolder, rootCounter)
		// since we have no idea of the generated meta file structure, we'll parse it into a map[string]interface{}
		var langExtraVars map[string]interface{}
		err := lib.Parse(rootMetaFile, &langExtraVars)
		if err != nil {
			log.Fatal(err)
		}
		// and then add the new language specific variables parsed from the meta file to the build matrix (which already contains the env variables)
		err = feedMatrix(langExtraVars, &mx)
		if err != nil {
			log.Fatal(err)
		}

		// we're not done yet: we need to also handle the matrix exclusions
		// we parse them into a list of matrices
		exclusions, err := exclusionsMatrices(config.Matrix.Exclude)
		if err != nil {
			log.Fatal(err)
		}

		// and finally, we iterate over all the permutations of the build matrix
		// these permutations are the list of all the combinations of env variables and language specific variables, minux the exclusions
		mx.IterAll(func(permutation map[string]string, counter string) {
			// we get called for every non-excluded permutations with the different variables values for this permutations and a unique permutation counter
			// handlePermutation will start from the .bazooka.*.yml file, which should already contain a single language specific permutation
			// and enrich it with the env variables combination
			// the same goes for the meta file
			if err := handlePermutation(permutation, config, counter, rootCounter); err != nil {
				log.Fatal(fmt.Errorf("Error while generating the permutations: %v", err))
			}
		}, exclusions)

		// after we're done iterating over the .bazooka.*.yml, and since we generated a new set of build files
		// we can now safely remove them
		err = os.Remove(file)
		if err != nil {
			log.Fatal(fmt.Errorf("Error while removing file %s: %v", file, err))
		}

		// same for the meta files
		err = os.Remove(fmt.Sprintf("%s/%s", MetaFolder, rootCounter))
		if err != nil {
			log.Fatal(fmt.Errorf("Error while removing meta folders: %v", err))
		}
	}
	log.Info("Matrix generated")

	log.Info("Starting generating Dockerfiles from Matrix")
	// Now we're left with the final build files
	files, err = lib.ListFilesWithPrefix(OutputFolder, ".bazooka")
	if err != nil {
		log.Fatal(fmt.Errorf("Error while listing .bazooka* files: %v", err))
	}

	for _, file := range files {
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			log.Fatal(fmt.Errorf("Error while parsing config file %s: %v", file, err))
		}

		// transform the .bazooka.x.yml file into a set of dockerfile + shell scripts who perform the actual build
		g := &Generator{
			Config:       config,
			OutputFolder: OutputFolder,
			Index:        parseCounter(file),
		}
		err = g.GenerateDockerfile()
		if err != nil {
			log.Fatal("Error while generating a dockerfile: %v", err)
		}
	}
	log.Info("Dockerfiles all created successfully")

}

func handlePermutation(permutation map[string]string, config *lib.Config, counter, rootCounter string) error {
	//Flush file
	// start from the language-spcecific permutation
	newConfig := *config

	// and replace its env variables with this unique permutation
	envMap := extractPrefixedKeysMap(permutation, MX_ENV_PREFIX)

	// Insert BZK_BUILD_DIR if not present
	if _, ok := envMap["BZK_BUILD_DIR"]; !ok {
		envMap["BZK_BUILD_DIR"] = "/bazooka"
	}

	newConfig.Env = lib.FlattenEnvMap(envMap)
	if err := lib.Flush(newConfig, fmt.Sprintf("%s/.bazooka.%s%s.yml", OutputFolder, rootCounter, counter)); err != nil {
		return err
	}

	// do the same for the meta file
	// start from the language specific permutation meta file
	rootMetaFile := fmt.Sprintf("%s/%s", MetaFolder, rootCounter)
	metaFile := fmt.Sprintf("%s/%s%s", MetaFolder, rootCounter, counter)

	// copy it to a global (lang specific+env vars) permutation meta file
	if err := lib.CopyFile(rootMetaFile, metaFile); err != nil {
		return err
	}
	// and add to it this unique permutation of env variables
	var buffer bytes.Buffer
	buffer.WriteString("env:\n")
	for _, env := range lib.FlattenEnvMap(envMap) {
		buffer.WriteString(fmt.Sprintf(" - %s\n", env))
	}
	// and write it to disk
	if err := lib.AppendToFile(metaFile, buffer.String(), 0644); err != nil {
		return err
	}

	return nil
}

func feedMatrix(extra map[string]interface{}, mx *matrix.Matrix) error {
	for k, v := range extra {
		switch k {
		case "env":
			if vs, ok := v.([]interface{}); ok {
				envVars := []string{}
				for _, envVar := range vs {
					if strEnvVar, ok := envVar.(string); ok {
						envVars = append(envVars, strEnvVar)
					} else {
						return fmt.Errorf("Invalid config: env should contain a sequence of strings: found a non string value %v:%T", envVar, envVar)

					}
				}
				mx.Merge(explodeProps(envVars, MX_ENV_PREFIX))
			} else {
				return fmt.Errorf("Invalid config: env should contain a sequence of strings: %v:%T", v, v)
			}

		default:
			mx.AddVar(k, fmt.Sprintf("%v", v))
		}
	}
	return nil
}

func exclusionsMatrices(xs []map[string]interface{}) ([]*matrix.Matrix, error) {
	res := make([]*matrix.Matrix, len(xs))
	for i, x := range xs {
		mx := matrix.Matrix{}
		if err := feedMatrix(x, &mx); err != nil {
			return nil, err
		}
		res[i] = &mx
	}
	return res, nil
}

// parseCounter extract the * part from a .bazooka.*.yml file name
func parseCounter(filePath string) string {
	splits := strings.Split(filePath, "/")
	file := splits[len(splits)-1]
	return strings.Split(file, ".")[2]
}

// explodeProps starts from a list of key=valye strings and stores them into a map
// it also handles repeated values, so ["A=1", "A=2", B="3"] gets transformed into {A: [1, 2], B: [3]}
func explodeProps(props []string, keyPrefix string) map[string][]string {
	envKeyMap := make(map[string][]string)
	for _, env := range props {
		envSplit := strings.Split(env, "=")
		envKeyMap[keyPrefix+envSplit[0]] = append(envKeyMap[keyPrefix+envSplit[0]], envSplit[1])
	}
	return envKeyMap
}

// prefixMapKeys returns a new map where all keys are prefixed with prefix
func prefixMapKeys(m map[string][]string, prefix string) map[string][]string {
	res := make(map[string][]string)
	for k, v := range m {
		res[prefix+k] = v
	}
	return res
}

// extractPrefixedKeysMap returns a new map containing only the values whose keys have the specified prefix, removing the latter in the process
// Given {xA: 1, B: 2, xC: 3}, it returns {A: 1, C: 3} if given a prefix x
func extractPrefixedKeysMap(m map[string]string, prefix string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		if strings.HasPrefix(k, prefix) {
			res[k[len(prefix):]] = v
		}
	}
	return res
}

func resolveLanguageParser(language string) (string, error) {
	parserMap := map[string]string{
		"golang": "bazooka/parser-golang",
		"go":     "bazooka/parser-golang",
		"java":   "bazooka/parser-java",
	}
	if val, ok := parserMap[language]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Language Parser %s\n", language)
}

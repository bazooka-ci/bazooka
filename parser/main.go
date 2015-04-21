package main

import (
	"fmt"
	"strings"

	"github.com/bazooka-ci/bazooka/commons/matrix"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	BazookaConfigFile = ".bazooka.yml"
	TravisConfigFile  = ".travis.yml"
	MX_ENV_PREFIX     = "env::"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
	err := lib.LoadCryptoKeyFromFile("/bazooka-cryptokey")
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	log.Info("Starting Parsing Phase")
	// Find either .travis.yml or .bazooka.yml file in the project
	configFile, err := lib.ResolveConfigFile(paths.container.source)
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

	var variants []*variantData

	if len(config.Language) == 0 {
		if len(config.Image) == 0 {
			log.Fatal("One of 'language' or 'image' needs to be set")
		}
		variants = generateImageVariants(config)
	} else {
		// resolve the docker image corresponding to this particular language parser
		image, err := resolveLanguageParser(config.Language)
		if err != nil {
			log.Fatal(err)
		}

		// run the parser image
		langParser := &LanguageParser{
			Image: image,
		}
		variants, err = langParser.Parse()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Starting Matrix generation")

	for _, variant := range variants {
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
		mx := matrix.Matrix(explodeProps(variant.config.Env, MX_ENV_PREFIX))

		// and then add the new language specific variables parsed from the meta file to the build matrix (which already contains the env variables)
		err = feedMatrix(variant.meta, &mx)
		if err != nil {
			log.Fatal(err)
		}

		// we're not done yet: we need to also handle the matrix exclusions
		// we parse them into a list of matrices
		exclusions, err := exclusionsMatrices(variant.config.Matrix.Exclude)
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
			if err := handlePermutation(permutation, variant.config, variant.meta, counter, variant.counter); err != nil {
				log.Fatal(fmt.Errorf("Error while generating the permutations: %v", err))
			}
		}, exclusions)
	}
	log.Info("Matrix generated")

	log.Info("Starting generating Dockerfiles from Matrix")
	// Now we're left with the final build files
	files, err := lib.ListFilesWithPrefix(paths.container.output, ".bazooka")
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
			OutputFolder: paths.container.output,
			Index:        parseCounter(file),
		}
		err = g.GenerateDockerfile()
		if err != nil {
			log.Fatal("Error while generating a dockerfile: %v", err)
		}
	}
	log.Info("Dockerfiles all created successfully")

}

func handlePermutation(permutation map[string]string, config *lib.Config, meta map[string]interface{}, counter, rootCounter string) error {
	// start from the language-spcecific permutation
	newConfig := *config

	// and replace its env variables with this unique permutation
	envMap := extractPrefixedKeysMap(permutation, MX_ENV_PREFIX)

	// Insert BZK_BUILD_DIR if not present
	if _, ok := envMap["BZK_BUILD_DIR"]; !ok {
		envMap["BZK_BUILD_DIR"] = "/bazooka"
	}

	newConfig.Env = lib.FlattenEnvMap(envMap)
	if err := lib.Flush(newConfig, fmt.Sprintf("%s/.bazooka.%s%s.yml", paths.container.output, rootCounter, counter)); err != nil {
		return err
	}

	// do the same for the meta file
	// start from the language specific permutation meta file
	// and add this permutation's env map
	meta["env"] = newConfig.Env
	metaFile := fmt.Sprintf("%s/%s%s", paths.container.meta, rootCounter, counter)
	lib.Flush(meta, metaFile)

	return nil
}

func generateImageVariants(conf *lib.Config) []*variantData {
	res := make([]*variantData, len(conf.Image))
	for i, im := range conf.Image {
		imageConf := *conf
		imageConf.FromImage = im
		imageConf.Image = nil

		res[i] = &variantData{
			counter: fmt.Sprintf("%d", i),
			config:  &imageConf,
			meta: map[string]interface{}{
				"image": im,
			},
		}
	}
	return res
}

func feedMatrix(extra map[string]interface{}, mx *matrix.Matrix) error {
	for k, v := range extra {
		switch k {
		case "env":
			if vs, ok := v.([]interface{}); ok {
				envVars := []lib.BzkString{}
				for _, envVar := range vs {
					if strEnvVar, ok := envVar.(string); ok {
						envVars = append(envVars, lib.BzkString(strEnvVar))
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
func explodeProps(props []lib.BzkString, keyPrefix string) map[string][]string {
	envKeyMap := make(map[string][]string)
	for _, env := range props {
		envSplit := strings.Split(string(env), "=")
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
func extractPrefixedKeysMap(m map[string]string, prefix string) map[lib.BzkString]lib.BzkString {
	res := make(map[lib.BzkString]lib.BzkString)
	for k, v := range m {
		if strings.HasPrefix(k, prefix) {
			res[lib.BzkString(k[len(prefix):])] = lib.BzkString(v)
		}
	}
	return res
}

func resolveLanguageParser(language string) (string, error) {
	parserMap := map[string]string{
		"golang":  "bazooka/parser-golang",
		"go":      "bazooka/parser-golang",
		"java":    "bazooka/parser-java",
		"python":  "bazooka/parser-python",
		"node_js": "bazooka/parser-nodejs",
		"nodejs":  "bazooka/parser-nodejs",
	}
	if val, ok := parserMap[language]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Language Parser %s\n", language)
}

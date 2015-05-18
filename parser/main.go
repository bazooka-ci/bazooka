package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bazooka-ci/bazooka/commons/matrix"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	BazookaConfigFile       = ".bazooka.yml"
	TravisConfigFile        = ".travis.yml"
	MX_ENV_PREFIX           = "env::"
	BazookaEnvJobParameters = "BZK_JOB_PARAMETERS"
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

	envParamString := os.Getenv(BazookaEnvJobParameters)
	var envParams []lib.BzkString
	err = json.Unmarshal([]byte(envParamString), &envParams)
	if err != nil {
		log.Fatal(err)
	}

	jobParameters := groupByName(envParams, MX_ENV_PREFIX)

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
		// groupByName transforms that into a map[string][]string
		// for example ["A=1", "A=2", B="3"]
		// when grouped by name gets transformed into
		// {A: [1, 2], B: [3]}
		// this matches the matrix layout, so we could store them directly into the matrix
		// but since we would like to be able to extract them later, and to avoid mixing them with a language specific variables (like jdk, go, etc.)
		// explode prefixes the env variables names with a prefix defined in the constant MX_ENV_PREFIX
		// Hence, our matrix is more like: {"env::A": [1, 2], "env::B": [3]}
		variantVariables := groupByName(variant.config.Env, MX_ENV_PREFIX)

		// insert or replace environment variables defined by the job parameters
		for k, v := range jobParameters {
			variantVariables[k] = v
		}

		// check if all environment variables are defined
		for k, v := range variantVariables {
			if len(v) == 1 && len(v[0]) == 0 {
				log.Fatalf("Missing required parameter %v", k)
			}
		}

		mx := matrix.Matrix(variantVariables)

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
	envMap := extractEnv(permutation, config.Env)

	// Insert BZK_BUILD_DIR if not present
	if _, ok := envMap["BZK_BUILD_DIR"]; !ok {
		envMap["BZK_BUILD_DIR"] = lib.BzkString{"BZK_BUILD_DIR", "/bazooka", false}
	}

	newConfig.Env = mapValues(envMap)
	if err := lib.Flush(newConfig, fmt.Sprintf("%s/.bazooka.%s%s.yml", paths.container.output, rootCounter, counter)); err != nil {
		return err
	}

	// do the same for the meta file
	// start from the language specific permutation meta file
	// and add this permutation's env map
	metaEnv, err:=generateEnvForMeta(newConfig.Env)
	if err != nil {
		return err
	}
	
	meta["env"] = metaEnv
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
			switch converted:=v.(type) {
			case []interface{}:
				envVars := []lib.BzkString{}
				for _, envVar := range converted {
					if strEnvVar, ok := envVar.(string); ok {
						n, v:=lib.SplitNameValue(strEnvVar)
						envVars = append(envVars, lib.BzkString{n, v, false})
					} else {
						return fmt.Errorf("Invalid config: env should contain a sequence of strings: found a non string value %v:%T", envVar, envVar)

					}
				}
				mx.Merge(groupByName(envVars, MX_ENV_PREFIX))
			case string:
						n, v:=lib.SplitNameValue(converted)
				mx.Merge(groupByName([]lib.BzkString{lib.BzkString{n, v, false}}, MX_ENV_PREFIX))
			default:
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

// groupByName starts from a list of key=valye strings and stores them into a map
// it also handles repeated values, so ["A=1", "A=2", B="3"] gets transformed into {A: [1, 2], B: [3]}
func groupByName(props []lib.BzkString, keyPrefix string) map[string][]string {
	res := make(map[string][]string)
	for _, env := range props {
		res[keyPrefix+env.Name] = append(res[keyPrefix+env.Name], env.Value)
	}
	return res
}

// prefixMapKeys returns a new map where all keys are prefixed with prefix
func prefixMapKeys(m map[string][]string, prefix string) map[string][]string {
	res := make(map[string][]string)
	for k, v := range m {
		res[prefix+k] = v
	}
	return res
}

// extractEnv extracts back the env variables from the permutation values
func extractEnv(from map[string]string, originalEnv []lib.BzkString) map[string]lib.BzkString {
	res := make(map[string]lib.BzkString)
	for k, v := range from {
		if strings.HasPrefix(k, MX_ENV_PREFIX) {
			name:=strings.TrimPrefix(k, MX_ENV_PREFIX)
			res[name] = findOrCreateBzkString(name, v, originalEnv)
		}
	}
	return res
}

func findOrCreateBzkString(name, value string, env []lib.BzkString) lib.BzkString {
	for _, e:=range env {
		if name==e.Name && value==e.Value {
			return e
		}
	}
	return lib.BzkString{name, value, false}
}

func mapValues(m map[string]lib.BzkString) []lib.BzkString {
	res :=make([]lib.BzkString, 0, len(m))
	for _, v:=range m {
		res=append(res, v)
	}
	return res
}

func generateEnvForMeta(env []lib.BzkString) ([]string, error) {
	key, err:=lib.ReadCryptoKey(paths.container.cryptoKey)
	if err != nil {
		return nil, err
	}
	
	res := make([]string, 0, len(env))
	for _, e:=range env {
		value := e.Value
		if e.Secured {
			eValue, err:= lib.Encrypt(key, []byte(value))
			if err != nil {
				return nil, err
			}
			value = hex.EncodeToString(eValue)
		}
		res = append(res, fmt.Sprintf("%s=%s", e.Name, value))
	}
	return res, nil
}
func resolveLanguageParser(language string) (string, error) {
	parserMap := map[string]string{
		"golang":  "bazooka/parser-golang",
		"go":      "bazooka/parser-golang",
		"java":    "bazooka/parser-java",
		"python":  "bazooka/parser-python",
		"node_js": "bazooka/parser-nodejs",
		"nodejs":  "bazooka/parser-nodejs",
		"php":     "bazooka/parser-php",
	}
	if val, ok := parserMap[language]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Language Parser %s\n", language)
}

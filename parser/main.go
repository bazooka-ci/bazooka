package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	lib "github.com/haklop/bazooka/commons"
)

const (
	SourceFolder      = "/bazooka"
	OutputFolder      = "/bazooka-output"
	BazookaConfigFile = ".bazooka.yml"
	TravisConfigFile  = ".travis.yml"
)

var permutationIndex int

func main() {

	configFile, err := lib.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}
	config := &lib.Config{}
	err = lib.Parse(configFile, config)
	if err != nil {
		log.Fatal(err)
	}
	image, err := resolveLanguageParser(config.Language)
	if err != nil {
		log.Fatal(err)
	}

	langParser := &LanguageParser{
		Image: image,
	}
	err = langParser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	files, err := lib.ListFilesWithPrefix(OutputFolder, ".bazooka")
	if err != nil {
		log.Fatal(err)
	}

	for i, file := range files {
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			log.Fatal(err)
		}

		permutations := permut(getEnvMap(config))
		err = iterPermutations(permutations, make(map[string]string), config, i)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Remove(file)
		if err != nil {
			log.Fatal(err)
		}
	}

	files, err = lib.ListFilesWithPrefix(OutputFolder, ".bazooka")
	if err != nil {
		log.Fatal(err)
	}

	for i, file := range files {
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			log.Fatal(err)
		}

		g := &Generator{
			Config:       config,
			OutputFolder: OutputFolder,
			Index:        i,
		}
		err = g.GenerateDockerfile()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func iterPermutations(perms []*Permutation, envMap map[string]string, config *lib.Config, rootIndex int) error {
	if len(perms) == 0 {
		//Flush file
		permutationIndex++
		newConfig := *config
		newConfig.Env = flattenMap(envMap)
		return lib.Flush(newConfig, fmt.Sprintf("%s/.bazooka.%d%d.yml", OutputFolder, rootIndex, permutationIndex))
	}
	for _, perm := range perms {
		envMap[perm.EnvKey] = perm.EnvValue
		iterPermutations(perm.Permutations, envMap, config, rootIndex)
	}
	return nil
}

func flattenMap(mapp map[string]string) []string {
	res := []string{}
	for key, value := range mapp {
		res = append(res, fmt.Sprintf("%s=%s", key, value))
	}
	return res
}

func permut(envKeyMap map[string][]string) []*Permutation {
	if len(envKeyMap) == 0 {
		return nil
	}
	var anyKey string
	for key := range envKeyMap {
		anyKey = key
		break
	}

	lowerMap := copyMap(envKeyMap)
	delete(lowerMap, anyKey)

	perms := []*Permutation{}
	for _, value := range envKeyMap[anyKey] {
		perms = append(perms, &Permutation{
			EnvKey:       anyKey,
			EnvValue:     value,
			Permutations: permut(lowerMap),
		})
	}
	return perms
}

type Permutation struct {
	EnvKey       string
	EnvValue     string
	Permutations []*Permutation
}

func copyMap(source map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for k, v := range source {
		dst[k] = v
	}
	return dst
}

func getEnvMap(config *lib.Config) map[string][]string {
	envKeyMap := make(map[string][]string)
	for _, env := range config.Env {
		envSplit := strings.Split(env, "=")
		envKeyMap[envSplit[0]] = append(envKeyMap[envSplit[0]], envSplit[1])
	}
	return envKeyMap
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

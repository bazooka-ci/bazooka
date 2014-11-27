package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	lib "github.com/haklop/bazooka/commons"
)

const (
	SourceFolder      = "/bazooka"
	OutputFolder      = "/bazooka-output"
	MetaFolder        = "/meta"
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

	for _, file := range files {
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			log.Fatal(err)
		}

		permutations := permut(getEnvMap(config))
		permutationIndex = 0
		err = iterPermutations(permutations, make(map[string]string), config, parseIndex(file))
		if err != nil {
			log.Fatal(err)
		}

		err = os.Remove(file)
		if err != nil {
			log.Fatal(err)
		}

		err = os.RemoveAll(fmt.Sprintf("%s/%s", MetaFolder, permutationIndex))
		if err != nil {
			log.Fatal(err)
		}
	}

	files, err = lib.ListFilesWithPrefix(OutputFolder, ".bazooka")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			log.Fatal(err)
		}

		g := &Generator{
			Config:       config,
			OutputFolder: OutputFolder,
			Index:        parseIndex(file),
		}
		err = g.GenerateDockerfile()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func parseIndex(filePath string) string {
	splits := strings.Split(filePath, "/")
	file := splits[len(splits)-1]
	return strings.Split(file, ".")[2]
}

func iterPermutations(perms []*Permutation, envMap map[string]string, config *lib.Config, rootIndex string) error {
	if len(perms) == 0 {
		//Flush file
		newConfig := *config

		if _, ok := envMap["BZK_BUILD_DIR"]; !ok {
			envMap["BZK_BUILD_DIR"] = "/bazooka"
		}

		newConfig.Env = lib.FlattenEnvMap(envMap)
		err := lib.CopyFile(fmt.Sprintf("%s/%s", MetaFolder, rootIndex), fmt.Sprintf("%s/%s%d", MetaFolder, rootIndex, permutationIndex))
		if err != nil {
			return err
		}
		var buffer bytes.Buffer
		buffer.WriteString("env:\n")
		for _, env := range lib.FlattenEnvMap(envMap) {
			buffer.WriteString(fmt.Sprintf(" - %s\n", env))
			if err != nil {
				return err
			}
		}
		err = lib.AppendToFile(fmt.Sprintf("%s/%s%d", MetaFolder, rootIndex, permutationIndex), buffer.String(), 0755)
		if err != nil {
			return err
		}
		err = lib.Flush(newConfig, fmt.Sprintf("%s/.bazooka.%s%d.yml", OutputFolder, rootIndex, permutationIndex))
		if err != nil {
			return err
		}

		permutationIndex++
	}
	for _, perm := range perms {
		envMap[perm.EnvKey] = perm.EnvValue
		iterPermutations(perm.Permutations, envMap, config, rootIndex)
	}
	return nil
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

	lowerMap := lib.CopyMap(envKeyMap)
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

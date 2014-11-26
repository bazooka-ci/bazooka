package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/haklop/bazooka/commons/matrix"

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
		rootCounter := parseCounter(file)
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			log.Fatal(fmt.Errorf("Error while parsing config file %s: %v", file, err))
		}
		mx := matrix.Matrix(getEnvMap(config))

		matrix.IterAll(mx, func(permutation map[string]string, counter string) {
			if err := handlePermutation(permutation, config, counter, rootCounter); err != nil {
				log.Fatal(fmt.Errorf("Error while generating the permutations: %v", err))
			}
		})

		err = os.Remove(file)
		if err != nil {
			log.Fatal(fmt.Errorf("Error while removing file %s: %v", file, err))
		}

		err = os.Remove(fmt.Sprintf("%s/%s", MetaFolder, rootCounter))
		if err != nil {
			log.Fatal(fmt.Errorf("Error while removing meta folders: %v", err))
		}
	}

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

		g := &Generator{
			Config:       config,
			OutputFolder: OutputFolder,
			Index:        parseCounter(file),
		}
		err = g.GenerateDockerfile()
		if err != nil {
			fmt.Errorf("Error while generating a dockerfile: %v", err)
		}
	}

}

func parseCounter(filePath string) string {
	splits := strings.Split(filePath, "/")
	file := splits[len(splits)-1]
	return strings.Split(file, ".")[2]
}

func handlePermutation(envMap map[string]string, config *lib.Config, counter, rootCounter string) error {

	//Flush file
	newConfig := *config
	newConfig.Env = lib.FlattenEnvMap(envMap)
	err := lib.CopyFile(fmt.Sprintf("%s/%s", MetaFolder, rootCounter), fmt.Sprintf("%s/%s%s", MetaFolder, rootCounter, counter))
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
	err = lib.AppendToFile(fmt.Sprintf("%s/%s%s", MetaFolder, rootCounter, counter), buffer.String(), 0755)
	if err != nil {
		return err
	}
	err = lib.Flush(newConfig, fmt.Sprintf("%s/.bazooka.%s%s.yml", OutputFolder, rootCounter, counter))
	if err != nil {
		return err
	}

	return nil
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

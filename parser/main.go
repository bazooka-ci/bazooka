package main

import (
	"fmt"
	"log"

	lib "github.com/bazooka-ci/bazooka-lib"
)

const (
	SourceFolder      = "/bazooka"
	OutputFolder      = "/bazooka-output"
	BazookaConfigFile = ".bazooka.yml"
	TravisConfigFile  = ".travis.yml"
)

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

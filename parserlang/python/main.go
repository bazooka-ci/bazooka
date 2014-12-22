package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/haklop/bazooka/commons/matrix"

	bazooka "github.com/haklop/bazooka/commons"
)

const (
	SourceFolder = "/bazooka"
	OutputFolder = "/bazooka-output"
	MetaFolder   = "/meta"
	PyLang       = "python"
)

type ConfigPython struct {
	Base       bazooka.Config `yaml:",inline"`
	PyVersions []string       `yaml:"python,omitempty"`
}

func main() {
	file, err := bazooka.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	conf := &ConfigPython{}
	err = bazooka.Parse(file, conf)
	if err != nil {
		log.Fatal(err)
	}

	if len(conf.Base.Script) == 0 {
		log.Fatal("Pyton builds should define a script value in the build descriptor")
	}

	mx := matrix.Matrix{
		PyLang: conf.PyVersions,
	}

	if len(conf.PyVersions) == 0 {
		mx[PyLang] = []string{"2.7"}
	}
	mx.IterAll(func(permutation map[string]string, counter string) {
		if err := managePyVersion(counter, conf, permutation[PyLang]); err != nil {
			log.Fatal(err)
		}
	}, nil)
}

func managePyVersion(counter string, conf *ConfigPython, version string) error {
	conf.PyVersions = []string{}
	image, err := resolvePyImage(version)
	conf.Base.FromImage = image
	if err != nil {
		return err
	}

	err = bazooka.AppendToFile(fmt.Sprintf("%s/%s", MetaFolder, counter), fmt.Sprintf("%s: %s\n", PyLang, version), 0644)
	if err != nil {
		return err
	}
	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func resolvePyImage(version string) (string, error) {
	//TODO extract this from db
	pyMap := map[string]string{
		"2.7": "bazooka/runner-python:2.7",
		"3.3": "bazooka/runner-python:3.3",
		"3.4": "bazooka/runner-python:3.4",
	}
	if val, ok := pyMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Python Runnner %s", version)
}

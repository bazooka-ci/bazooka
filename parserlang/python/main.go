package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	bazooka "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	SourceFolder = "/bazooka"
	OutputFolder = "/bazooka-output"
	MetaFolder   = "/meta"
	PyLang       = "python"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

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

	versions := conf.PyVersions
	images := conf.Base.Image

	if len(versions) == 0 && len(images) == 0 {
		versions = []string{"2.7"}
	}

	for i, version := range versions {
		if err := managePyVersion(fmt.Sprintf("0%d", i), conf, version, ""); err != nil {
			log.Fatal(err)
		}
	}
	for i, image := range images {
		if err := managePyVersion(fmt.Sprintf("1%d", i), conf, "", image); err != nil {
			log.Fatal(err)
		}
	}
}

func managePyVersion(counter string, conf *ConfigPython, version, image string) error {
	conf.PyVersions = nil
	conf.Base.Image = nil

	meta := map[string]string{}
	if len(version) > 0 {
		var err error
		image, err = resolvePyImage(version)
		if err != nil {
			return err
		}
		meta[PyLang] = version
	} else {
		meta["image"] = image
	}
	conf.Base.FromImage = image

	if err := bazooka.Flush(meta, fmt.Sprintf("%s/%s", MetaFolder, counter)); err != nil {
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

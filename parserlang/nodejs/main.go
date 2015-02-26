package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bazooka-ci/bazooka/commons/matrix"

	bazooka "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	SourceFolder = "/bazooka"
	OutputFolder = "/bazooka-output"
	MetaFolder   = "/meta"
	Nodejs       = "node_js"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

type ConfigNodejs struct {
	Base         bazooka.Config `yaml:",inline"`
	NodeVersions []string       `yaml:"node_js,omitempty"`
}

func main() {
	file, err := bazooka.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	conf := &ConfigNodejs{}
	err = bazooka.Parse(file, conf)
	if err != nil {
		log.Fatal(err)
	}

	mx := matrix.Matrix{
		Nodejs: conf.NodeVersions,
	}

	if len(conf.NodeVersions) == 0 {
		mx[Nodejs] = []string{"0.10"}
	}
	mx.IterAll(func(permutation map[string]string, counter string) {
		if err := manageNodejsVersion(counter, conf, permutation[Nodejs]); err != nil {
			log.Fatal(err)
		}
	}, nil)
}

func manageNodejsVersion(counter string, conf *ConfigNodejs, version string) error {
	conf.NodeVersions = []string{}
	image, err := resolveNodejsImage(version)
	if err != nil {
		return err
	}
	conf.Base.FromImage = image

	setDefaultInstall(conf)
	setDefaultScript(conf)

	err = bazooka.AppendToFile(fmt.Sprintf("%s/%s", MetaFolder, counter), fmt.Sprintf("%s: %s\n", Nodejs, version), 0644)
	if err != nil {
		return err
	}
	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func resolveNodejsImage(version string) (string, error) {
	//TODO extract this from db
	nodeMap := map[string]string{
		"0.8":  "bazooka/runner-nodejs:0.8",
		"0.10": "bazooka/runner-nodejs:0.10",
		"0.11": "bazooka/runner-nodejs:0.11",
	}
	if val, ok := nodeMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for NodeJS Runnner %s", version)
}

func setDefaultInstall(conf *ConfigNodejs) {
	if len(conf.Base.Install) == 0 {
		conf.Base.Install = []string{"npm install"}
	}
}

func setDefaultScript(conf *ConfigNodejs) {
	if len(conf.Base.Script) == 0 {
		conf.Base.Script = []string{"npm test"}
	}
}

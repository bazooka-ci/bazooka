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

	versions := conf.NodeVersions
	images := conf.Base.Image

	if len(versions) == 0 && len(images) == 0 {
		versions = []string{"0.10"}
	}
	for i, version := range versions {
		if err := manageNodejsVersion(fmt.Sprintf("0%d", i), conf, version, ""); err != nil {
			log.Fatal(err)
		}
	}
	for i, image := range images {
		if err := manageNodejsVersion(fmt.Sprintf("1%d", i), conf, "", image); err != nil {
			log.Fatal(err)
		}
	}

}

func manageNodejsVersion(counter string, conf *ConfigNodejs, version, image string) error {
	conf.NodeVersions = nil
	conf.Base.Image = nil

	setDefaultInstall(conf)
	setDefaultScript(conf)

	meta := map[string]string{}
	if len(version) > 0 {
		var err error
		image, err = resolveNodejsImage(version)
		if err != nil {
			return err
		}
		meta[Nodejs] = version
	} else {
		meta["image"] = image
	}
	conf.Base.FromImage = image

	if err := bazooka.Flush(meta, fmt.Sprintf("%s/%s", MetaFolder, counter)); err != nil {
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

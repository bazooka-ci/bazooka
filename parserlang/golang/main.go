package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"

	bazooka "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

const (
	SourceFolder = "/bazooka"
	OutputFolder = "/bazooka-output"
	MetaFolder   = "/meta"
	Golang       = "go"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

func main() {
	file, err := bazooka.ResolveConfigFile(SourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	conf := &ConfigGolang{}
	err = bazooka.Parse(file, conf)
	if err != nil {
		log.Fatal(err)
	}

	versions := conf.GoVersions
	images := conf.Base.Image

	if len(versions) == 0 && len(images) == 0 {
		versions = []string{"tip"}
	}
	for i, version := range versions {
		if err := manageGoVersion(fmt.Sprintf("0%d", i), conf, version, ""); err != nil {
			log.Fatal(err)
		}
	}

	for i, image := range images {
		if err := manageGoVersion(fmt.Sprintf("1%d", i), conf, "", image); err != nil {
			log.Fatal(err)
		}
	}
}

func manageGoVersion(counter string, conf *ConfigGolang, version string, image string) error {
	conf.GoVersions = nil
	conf.Base.Image = nil

	setGodir(conf)
	setDefaultInstall(conf)
	err := setDefaultScript(conf)
	if err != nil {
		return err
	}

	meta := map[string]string{}
	if len(version) > 0 {
		image, err = resolveGoImage(version)
		if err != nil {
			return err
		}
		meta[Golang] = version
	} else {
		meta["image"] = image
	}
	conf.Base.FromImage = image

	if err = bazooka.Flush(meta, fmt.Sprintf("%s/%s", MetaFolder, counter)); err != nil {
		return err
	}
	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func setGodir(conf *ConfigGolang) {
	env := bazooka.GetEnvMap(conf.Base.Env)

	godirExist, err := bazooka.FileExists("/bazooka/.godir")
	if err != nil {
		log.Fatal(err)
	}

	var buildDir string
	if godirExist {
		f, err := os.Open("/bazooka/.godir")
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}

		bf := bufio.NewReader(f)

		// only read first line
		content, isPrefix, err := bf.ReadLine()

		if err == io.EOF {
			buildDir = "/go/src/app"
		} else if err != nil {
			log.Fatal(err)
		} else if isPrefix {
			log.Fatal("Unexpected long line reading", f.Name())
		} else {
			buildDir = fmt.Sprintf("/go/src/%s", content)
		}

	} else {
		scmMetadata := &bazooka.SCMMetadata{}
		scmMetadataFile := fmt.Sprintf("%s/scm", MetaFolder)
		err = bazooka.Parse(scmMetadataFile, scmMetadata)
		if err != nil {
			log.Fatal(err)
		}

		if len(scmMetadata.Origin) > 0 {
			r, err := regexp.Compile("^(?:https://(?:\\w+@){0,1}|git@)(github.com|bitbucket.org)[:/]{0,1}([\\w-_]+/[\\w-_]+).git$")
			if err != nil {
				log.Fatal(err)
			}

			res := r.FindStringSubmatch(scmMetadata.Origin)
			if res != nil {
				buildDir = fmt.Sprintf("/go/src/%s/%s", res[1], res[2])
			} else {
				buildDir = "/go/src/app"
			}
		} else {
			buildDir = "/go/src/app"
		}
	}

	log.WithFields(log.Fields{
		"build_directory": buildDir,
	}).Info("Build directory set")

	env["BZK_BUILD_DIR"] = []string{buildDir}

	conf.Base.Env = flattenEnvMap(env)
}

func setDefaultInstall(conf *ConfigGolang) {
	if len(conf.Base.Install) == 0 {
		conf.Base.Install = []string{"go get -d -t -v ./... && go build -v ./..."}
	}
}

func setDefaultScript(conf *ConfigGolang) error {
	if len(conf.Base.Script) == 0 {
		if _, err := os.Open(fmt.Sprintf("%s/Makefile", SourceFolder)); err != nil {
			if os.IsNotExist(err) {
				conf.Base.Script = []string{"go test -v ./..."}
				return nil
			}
			return err
		}
		conf.Base.Script = []string{"make"}
	}
	return nil
}

func resolveGoImage(version string) (string, error) {
	//TODO extract this from db
	goMap := map[string]string{
		"1.2.2":  "bazooka/runner-golang:1.2.2",
		"1.3":    "bazooka/runner-golang:1.3",
		"1.3.1":  "bazooka/runner-golang:1.3.1",
		"1.3.2":  "bazooka/runner-golang:1.3.2",
		"1.3.3":  "bazooka/runner-golang:1.3.3",
		"1.4":    "bazooka/runner-golang:1.4",
		"1.4.1":  "bazooka/runner-golang:1.4.1",
		"1.4.2":  "bazooka/runner-golang:1.4.2",
		"tip":    "bazooka/runner-golang:latest",
		"latest": "bazooka/runner-golang:latest",
	}
	if val, ok := goMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Go Runnner %s\n", version)
}

func flattenEnvMap(mapp map[string][]string) []string {
	res := []string{}
	for key, values := range mapp {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return res
}

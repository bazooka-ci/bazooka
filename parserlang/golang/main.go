package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	bazooka "github.com/haklop/bazooka/commons"
)

const (
	SourceFolder      = "/bazooka"
	OutputFolder      = "/bazooka-output"
	MetaFolder        = "/meta"
	BazookaConfigFile = ".bazooka.yml"
	TravisConfigFile  = ".travis.yml"
)

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

	if len(conf.GoVersions) == 0 {
		err = manageGoVersion(0, conf, "tip")
		if err != nil {
			log.Fatal(err)
		}

	} else {
		for i, version := range conf.GoVersions {
			vconf := *conf
			err = manageGoVersion(i, &vconf, version)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func manageGoVersion(i int, conf *ConfigGolang, version string) error {
	conf.GoVersions = []string{}
	setGodir(conf)
	setDefaultInstall(conf)
	err := setDefaultScript(conf)
	if err != nil {
		return err
	}
	image, err := resolveGoImage(version)
	conf.FromImage = image
	if err != nil {
		return err
	}

	err = bazooka.AppendToFile(fmt.Sprintf("%s/%d", MetaFolder, i), fmt.Sprintf("golang: %s\n", version), 0644)
	if err != nil {
		return err
	}
	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%d.yml", OutputFolder, i))
}

func setGodir(conf *ConfigGolang) {
	env := bazooka.GetEnvMap(conf.Env)

	godirExist, err := bazooka.FileExists("/bazooka/.godir")
	if err != nil {
		log.Fatal(err)
	}

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
			env["BZK_BUILD_DIR"] = "/go/src/app"
		} else if err != nil {
			log.Fatal(err)
		} else if isPrefix {
			log.Fatal("Unexpected long line reading", f.Name())
		} else {
			env["BZK_BUILD_DIR"] = fmt.Sprintf("/go/src/%s", content)
		}

	} else {
		env["BZK_BUILD_DIR"] = "/go/src/app"
	}

	conf.Env = bazooka.FlattenEnvMap(env)
}

func setDefaultInstall(conf *ConfigGolang) {
	if len(conf.Install) == 0 {
		conf.Install = []string{"go get -d -v ./... && go build -v ./..."}
	}
}

func setDefaultScript(conf *ConfigGolang) error {
	if len(conf.Script) == 0 {
		if _, err := os.Open(fmt.Sprintf("%s/Makefile", SourceFolder)); err != nil {
			if os.IsNotExist(err) {
				conf.Script = []string{"go test -v ./..."}
				return nil
			}
			return err
		}
		conf.Script = []string{"make"}
	}
	return nil
}

func resolveGoImage(version string) (string, error) {
	//TODO extract this from db
	goMap := map[string]string{
		"1.2.2": "bazooka/runner-golang:1.2.2",
		"1.3":   "bazooka/runner-golang:1.3",
		"1.3.1": "bazooka/runner-golang:1.3.1",
		"1.3.2": "bazooka/runner-golang:1.3.2",
		"1.3.3": "bazooka/runner-golang:1.3.3",
		"tip":   "bazooka/runner-golang:latest",
	}
	if val, ok := goMap[version]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for Go Runnner %s\n", version)
}

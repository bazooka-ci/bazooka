package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/haklop/bazooka/commons/matrix"

	bazooka "github.com/haklop/bazooka/commons"
)

const (
	SourceFolder      = "/bazooka"
	OutputFolder      = "/bazooka-output"
	MetaFolder        = "/meta"
	BazookaConfigFile = ".bazooka.yml"
	TravisConfigFile  = ".travis.yml"
	Golang            = "go"
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

	mx := matrix.Matrix{
		Golang: conf.GoVersions,
	}

	if len(conf.GoVersions) == 0 {
		mx[Golang] = []string{"tip"}
	}
	mx.IterAll(func(permutation map[string]string, counter string) {
		if err := manageGoVersion(counter, conf, permutation[Golang]); err != nil {
			log.Fatal(err)
		}
	}, nil)
}

func manageGoVersion(counter string, conf *ConfigGolang, version string) error {
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

	err = bazooka.AppendToFile(fmt.Sprintf("%s/%s", MetaFolder, counter), fmt.Sprintf("%s: %s\n", Golang, version), 0644)
	if err != nil {
		return err
	}
	return bazooka.Flush(conf, fmt.Sprintf("%s/.bazooka.%s.yml", OutputFolder, counter))
}

func setGodir(conf *ConfigGolang) {
	fmt.Println("Conf: %#v", conf)
	env := bazooka.GetEnvMap(conf.Env)

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
		buildDir = "/go/src/app"
	}

	env["BZK_BUILD_DIR"] = []string{buildDir}


	conf.Env = flattenEnvMap(env)
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

func flattenEnvMap(mapp map[string][]string) []string {
	res := []string{}
		for key, values := range mapp {
			for _, value := range values {
				res = append(res, fmt.Sprintf("%s=%s", key, value))
			}
		}
	return res
}

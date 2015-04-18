package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	docker "github.com/bywan/go-dockercommand"
)

const (
	dockerEndpoint = "unix:///docker.sock"
)

type LanguageParser struct {
	Image string
}

type variantData struct {
	counter string
	config  *lib.Config
	meta    map[string]interface{}
}

func (p *LanguageParser) Parse() ([]*variantData, error) {

	log.WithFields(log.Fields{
		"image": p.Image,
	}).Info("Lauching language parsing")

	client, err := docker.NewDocker(dockerEndpoint)
	if err != nil {
		return nil, err
	}
	
	container, err := client.Run(&docker.RunOptions{
		Image: p.Image,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", paths.host.source),
			fmt.Sprintf("%s:/bazooka-output", paths.host.output),
			fmt.Sprintf("%s:/meta", paths.host.meta),
		},
		Detach: true,
	})
	if err != nil {
		return nil, err
	}

	container.Logs(p.Image)

	exitCode, err := container.Wait()
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, fmt.Errorf("Error during execution of Language Parser container %s/parser\n Check Docker container logs, id is %s\n", p.Image, container.ID())
	}

	err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil {
		return nil, err
	}
	log.Info("Language parsing finished")

	// if all went well, the parser should have generated one or more "sub" .bazooka.*.yml files
	// one for each compiler version for example
	//
	// they are also supposed to enrich it with a from attribute corresponding to a base docker image
	// to be used to run the build

	files, err := lib.ListFilesWithPrefix(paths.container.output, ".bazooka")
	if err != nil {
		log.Fatal(err)
	}

	res := make([]*variantData, len(files))

	// for each of those files (the "sub" .bazooka.*.yml)
	for i, file := range files {
		// parse the damned thing
		config := &lib.Config{}
		err = lib.Parse(file, config)
		if err != nil {
			return nil, err
		}

		// extract the "*" part from the .bazooka.*.yml file
		rootCounter := parseCounter(file)
		// for every .bazooka.*.yml file, the language parser is also supposed to have generated a meta/* file
		// which is a simple yml file containing the language specific  matrix variables
		// for example, if the original .bazooka.yml file defined 2 go versions:
		//
		// go:
		// - 1.2.2
		// - 1.3.1
		//
		// the language parser should generate 2 meta files, one for each go version in this format:
		//
		// go: 1.2.2
		//
		// and
		//
		// go: 1.3.1
		rootMetaFile := fmt.Sprintf("%s/%s", paths.container.meta, rootCounter)
		// since we have no idea of the generated meta file structure, we'll parse it into a map[string]interface{}
		var langExtraVars map[string]interface{}
		err := lib.Parse(rootMetaFile, &langExtraVars)
		if err != nil {
			return nil, err
		}
		res[i] = &variantData{
			counter: rootCounter,
			config:  config,
			meta:    langExtraVars,
		}
		// after we're done iterating over the .bazooka.*.yml, and since we generated a new set of build files
		// we can now safely remove them
		err = os.Remove(file)
		if err != nil {
			return nil, fmt.Errorf("Error while removing file %s: %v", file, err)
		}

		// same for the meta files
		err = os.Remove(rootMetaFile)
		if err != nil {
			return nil, fmt.Errorf("Error while removing meta folders: %v", err)
		}
	}
	return res, nil
}

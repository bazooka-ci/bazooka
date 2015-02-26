package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	lib "github.com/bazooka-ci/bazooka/commons"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
	"github.com/bazooka-ci/bazooka/commons/mongo"
)

type Parser struct {
	MongoConnector *mongo.MongoConnector
	Options        *ParseOptions
}

type ParseOptions struct {
	InputFolder    string
	OutputFolder   string
	DockerSock     string
	HostBaseFolder string
	MetaFolder     string
	Env            map[string]string
}

type variantData struct {
	counter    string
	meta       *lib.VariantMetas
	dockerFile string
	scripts    []string
	variant    *lib.Variant
	imageTag   string
}

func (p *Parser) Parse(logger Logger) ([]*variantData, error) {
	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return nil, err
	}

	image, err := p.resolveImage()
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"image": image,
	}).Info("Running Parsing Image on checked-out source")

	container, err := client.Run(&docker.RunOptions{
		Image: image,
		Env:   p.Options.Env,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", p.Options.InputFolder),
			fmt.Sprintf("%s:/meta", p.Options.MetaFolder),
			fmt.Sprintf("%s:/bazooka-output", p.Options.OutputFolder),
			fmt.Sprintf("%s:/docker.sock", p.Options.DockerSock)},
		Detach: true,
	})
	if err != nil {
		return nil, err
	}

	container.Logs(image)
	logger(image, "", container)

	exitCode, err := container.Wait()
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, fmt.Errorf("Error during execution of Parser container %s/parser\n Check Docker container logs, id is %s\n", image, container.ID())
	}

	err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"dockerfiles_path": p.Options.OutputFolder,
	}).Info("Parsing Image ran sucessfully, Dockerfiles generated")

	return p.variantsData()
}

func (p *Parser) variantsData() ([]*variantData, error) {
	workFolder := "/bazooka/work"
	metaFolder := "/bazooka/meta"
	dirs, err := ioutil.ReadDir(workFolder)
	if err != nil {
		return nil, fmt.Errorf("Failed to list the output dir files (%s): %v", workFolder, err)
	}
	var output []*variantData
	for _, dir := range dirs {
		if dir.Mode().IsDir() {
			vf := &variantData{counter: dir.Name()}
			files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", workFolder, dir.Name()))
			if err != nil {
				return nil, fmt.Errorf("Failed to list the variant %s files: %v", dir.Name(), err)
			}
			for _, file := range files {
				fullName := fmt.Sprintf("%s/%s/%s", workFolder, dir.Name(), file.Name())
				if file.Name() == "Dockerfile" {
					vf.dockerFile = fullName
					continue
				}
				vf.scripts = append(vf.scripts, fullName)
			}
			if len(vf.dockerFile) == 0 {
				return nil, fmt.Errorf("The variant %s has no Dockerfile", dir.Name())
			}
			if len(vf.scripts) == 0 {
				return nil, fmt.Errorf("The variant %s has no scripts", dir.Name())
			}

			metaFile := fmt.Sprintf("%s/%s", metaFolder, dir.Name())

			if err := parseMeta(metaFile, vf); err != nil {
				return nil, fmt.Errorf("Error while parsing that variant %s meta: %v", vf.counter, err)
			}
			sort.Sort(vf.meta)
			output = append(output, vf)
		}
	}
	return output, nil
}

func (f *Parser) resolveImage() (string, error) {
	image, err := f.MongoConnector.GetImage("parser")
	if err != nil {
		return "", fmt.Errorf("Unable to find Bazooka Docker Image for parser\n")
	}
	return image, nil
}

func parseMeta(file string, vf *variantData) error {
	meta := map[string]interface{}{}
	if err := lib.Parse(file, &meta); err != nil {
		return err
	}
	vf.meta = &lib.VariantMetas{}

	for k, v := range meta {
		switch k {
		case "env":
			if vs, ok := v.([]interface{}); ok {
				for _, envVar := range vs {
					if strEnvVar, ok := envVar.(string); ok {
						if strings.HasPrefix(strEnvVar, "BZK_") {
							continue
						}
						nv := strings.Split(strEnvVar, "=")
						vf.meta.Append(&lib.VariantMeta{Kind: lib.META_ENV, Name: nv[0], Value: nv[1]})
					} else {
						return fmt.Errorf("Invalid config: env should contain a sequence of strings: found a non string value %v:%T", envVar, envVar)

					}
				}
			} else {
				return fmt.Errorf("Invalid config: env should contain a sequence of strings: %v:%T", v, v)
			}

		default:
			vf.meta.Append(&lib.VariantMeta{Kind: lib.META_LANG, Name: k, Value: fmt.Sprintf("%v", v)})

		}
	}
	return nil
}

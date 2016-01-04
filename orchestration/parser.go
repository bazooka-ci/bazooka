package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	lib "github.com/bazooka-ci/bazooka/commons"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
)

type Parser struct {
	context *context
}

type variantData struct {
	counter    string
	meta       *lib.VariantMetas
	dockerFile string
	scripts    []string
	variant    *lib.Variant
	imageTag   string
	services   []lib.Service
}

func (p *Parser) Parse() ([]*variantData, error) {
	paths := p.context.paths
	client, err := docker.NewDocker(paths.dockerEndpoint.container)
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

	env := map[string]string{
		BazookaEnvApiUrl:        p.context.apiUrl,
		BazookaEnvSyslogUrl:     p.context.syslogUrl,
		BazookaEnvNetwork:       p.context.network,
		BazookaEnvHome:          paths.base.host,
		BazookaEnvSrc:           paths.source.host,
		BazookaEnvProjectID:     p.context.projectID,
		BazookaEnvJobID:         p.context.jobID,
		BazookaEnvJobParameters: p.context.jobParameters,
	}

	volumes := []string{
		fmt.Sprintf("%s:/bazooka", paths.source.host),
		fmt.Sprintf("%s:/meta", paths.meta.host),
		fmt.Sprintf("%s:/bazooka-output", paths.work.host),
		fmt.Sprintf("%s:/var/run/docker.sock", paths.dockerSock.host),
	}

	if len(paths.cryptoKey.host) > 0 {
		volumes = append(volumes, fmt.Sprintf("%s:/bazooka-cryptokey", paths.cryptoKey.host))
		env[BazookaEnvCryptoKeyfile] = paths.cryptoKey.host
	}

	container, err := client.Run(&docker.RunOptions{
		Image:               image,
		Env:                 env,
		VolumeBinds:         volumes,
		Detach:              true,
		NetworkMode:         p.context.network,
		LoggingDriver:       "syslog",
		LoggingDriverConfig: p.context.loggerConfig(image, ""),
	})
	if err != nil {
		return nil, err
	}

	defer lib.RemoveContainer(container)

	exitCode, err := container.Wait()
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, fmt.Errorf("Error during execution of Parser container %s/parser\n Check Docker container logs, id is %s\n", image, container.ID())
	}

	log.WithFields(log.Fields{
		"dockerfiles_path": paths.work.host,
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
				switch file.Name() {
				case "Dockerfile":
					vf.dockerFile = fullName
				case "services":
					servicesList := []lib.Service{}
					err := lib.Parse(fullName, &servicesList)
					if err != nil {
						return nil, fmt.Errorf("Failed to parse services file %s: %v", fullName, err)
					}
					vf.services = servicesList
				default:
					vf.scripts = append(vf.scripts, fullName)
				}
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
	image, err := f.context.client.Image.Get("parser")
	if err != nil {
		return "", fmt.Errorf("Unable to find Bazooka Docker Image for parser\n")
	}
	return image.Image, nil
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
						nv := strings.SplitN(strEnvVar, "=", 2)
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

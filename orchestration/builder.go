package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
	lib "github.com/haklop/bazooka/commons"
)

type Builder struct {
	Options *BuildOptions
}

type BuildOptions struct {
	DockerfileFolder string
	SourceFolder     string
	ProjectID        string
	JobID            string
}

type BuiltImage struct {
	Image     string
	VariantID int
}

func (b *Builder) Build() ([]BuiltImage, error) {

	log.Info("Starting building Dockerfiles")
	files, err := listBuildfiles(b.Options.DockerfileFolder)
	if err != nil {
		return nil, err
	}

	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	successChan := make(chan BuiltImage)
	remainingBuilds := len(files)

	for i, file := range files {
		go buildContainer(client, i, b, file, successChan, errChan)
	}

	var buildImages []BuiltImage
	for {
		select {
		case tag := <-successChan:
			buildImages = append(buildImages, tag)
			remainingBuilds--
		case err := <-errChan:
			return nil, err
		}

		if remainingBuilds == 0 {
			break
		}
	}
	return buildImages, nil
}

func buildContainer(client *docker.Docker, variantID int, b *Builder, file *buildFiles, successChan chan BuiltImage, errChan chan error) {
	for _, buildFile := range file.BuildFiles {
		splitString := strings.Split(buildFile, "/")
		err := lib.CopyFile(buildFile, fmt.Sprintf("%s/%s", b.Options.SourceFolder, splitString[len(splitString)-1]))
		if err != nil {
			errChan <- err
			return
		}
	}

	tag := fmt.Sprintf("bazooka/build-%s-%s-%d", b.Options.ProjectID, b.Options.JobID, variantID)
	err := client.Build(&docker.BuildOptions{
		Tag:        tag,
		Dockerfile: file.Dockerfile,
		Path:       b.Options.SourceFolder,
	})
	if err != nil {
		errChan <- err
	} else {
		successChan <- BuiltImage{
			Image:     tag,
			VariantID: variantID,
		}
	}
}

func listBuildfiles(source string) ([]*buildFiles, error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return nil, err
	}
	var output []*buildFiles
	for _, file := range files {
		if file.Mode().IsDir() {
			index, err := strconv.ParseInt(file.Name(), 10, 64)
			if err != nil {
				return nil, err
			}
			filesBuild, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", source, file.Name()))
			if err != nil {
				return nil, err
			}
			var result []string
			for _, fileBuild := range filesBuild {
				if fileBuild.Name() != "Dockerfile" {
					result = append(result, fmt.Sprintf("%s/%s/%s", source, file.Name(), fileBuild.Name()))
				}
			}
			output = append(output, &buildFiles{
				Dockerfile: fmt.Sprintf("%s/%s/Dockerfile", source, file.Name()),
				BuildFiles: result,
				JobIndex:   index,
			})

		}
	}
	return output, nil
}

type buildFiles struct {
	Dockerfile string
	BuildFiles []string
	JobIndex   int64
}

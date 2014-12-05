package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/haklop/bazooka/commons/parallel"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
	lib "github.com/haklop/bazooka/commons"
)

type Builder struct {
	Options *BuildOptions
}

type BuildOptions struct {
	SourceFolder string
	ProjectID    string
	Variants     []*variantData
}

func (b *Builder) Build() error {
	log.Info("Starting building Dockerfiles")

	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return err
	}

	par := parallel.New()
	for _, ivariant := range b.Options.Variants {
		variant := ivariant
		par.Submit(func() error {
			return b.buildContainer(client, variant)
		}, variant)
	}

	par.Exec(func(tag interface{}, err error) {
		v := tag.(*variantData)
		if err != nil {
			log.Errorf("Build error %v for variant %v\n", err, v)
			v.variant.Status = lib.JOB_ERRORED
			v.variant.Completed = time.Now()
			return
		}
		log.Infof("Build success for variant %v\n", v)
	})
	return nil
}

func (b *Builder) buildContainer(client *docker.Docker, vd *variantData) error {
	log.Infof("build container for variant %#v\n", vd)

	for _, script := range vd.scripts {

		splitString := strings.Split(script, "/")
		dest := fmt.Sprintf("%s/%s", b.Options.SourceFolder, splitString[len(splitString)-1])
		err := lib.CopyFile(script, dest)
		if err != nil {
			return err
		}
	}

	tag := fmt.Sprintf("bazooka/build-%s-%s-%d", b.Options.ProjectID, vd.variant.JobID, vd.variant.Number)

	err := client.Build(&docker.BuildOptions{
		Tag:        tag,
		Dockerfile: vd.dockerFile,
		Path:       b.Options.SourceFolder,
	})
	if err != nil {
		return err
	}
	vd.imageTag = tag
	return nil
}

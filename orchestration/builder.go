package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bazooka-ci/bazooka/commons/parallel"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	docker "github.com/bywan/go-dockercommand"
)

type Builder struct {
	context  *context
	variants []*variantData
}

func (b *Builder) Build() error {
	log.Info("Starting building Dockerfiles")
	paths := b.context.paths

	client, err := docker.NewDocker(paths.dockerEndpoint.container)
	if err != nil {
		return err
	}

	par := parallel.New()
	for _, ivariant := range b.variants {
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
		log.WithFields(log.Fields{
			"variant": v.counter,
		}).Info("Build success for variant")

	})
	return nil
}

func (b *Builder) buildContainer(client *docker.Docker, vd *variantData) error {
	log.WithFields(log.Fields{
		"variant": vd.counter,
	}).Info("Building container for variant")

	tag := fmt.Sprintf("bazooka-build/%s-%s-%d", b.context.projectID, vd.variant.JobID, vd.variant.Number)

	err := client.Build(&docker.BuildOptions{
		Tag:        tag,
		Dockerfile: strings.TrimPrefix(vd.dockerFile, "/bazooka/"), //Ugly hack: the Dockerfile path needs to be relative to context dir
		Path:       b.context.paths.base.container,
	})
	if err != nil {
		return err
	}
	vd.imageTag = tag
	return nil
}

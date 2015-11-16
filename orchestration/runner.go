package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	commons "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/commons/parallel"
	docker "github.com/bywan/go-dockercommand"
)

type Runner struct {
	variants []*variantData
	context  *context
	client   *docker.Docker
}

func (r *Runner) Run() error {
	paths := r.context.paths

	client, err := docker.NewDocker(paths.dockerEndpoint.container)
	if err != nil {
		return err
	}
	r.client = client

	par := parallel.New()

	for _, ivariant := range r.variants {
		if ivariant.variant.Status != commons.JOB_RUNNING {
			continue
		}
		variant := ivariant
		par.Submit(func() error {
			return r.runContainer(variant)
		}, variant)
	}

	par.Exec(func(tag interface{}, err error) {
		v := tag.(*variantData)
		if err != nil {
			log.Errorf("Run error %v for variant %v\n", err, v)
			v.variant.Status = commons.JOB_ERRORED
		} else {
			log.WithFields(log.Fields{
				"variant": v.counter,
			}).Info("Variant Completed")
		}
		v.variant.Completed = time.Now()
	})

	log.Info("Dockerfiles builds finished")
	return nil
}

func (r *Runner) runContainer(vd *variantData) error {
	paths := r.context.paths

	success := true

	serviceContainers := []*docker.Container{}
	containerLinks := []string{}
	for sidx, service := range vd.services {
		name := fmt.Sprintf("bazooka-service-%s-%s-%d-%d", r.context.projectID, r.context.jobID, vd.variant.Number, sidx)
		if len(service.Alias) == 0 {
			service.Alias = safeDockerAlias(strings.Split(service.Image, ":")[0])
		}
		containerLinks = append(containerLinks, fmt.Sprintf("%s:%s", name, service.Alias))
		serviceContainer, err := r.client.Run(&docker.RunOptions{
			Name:   name,
			Image:  service.Image,
			Detach: true,
		})
		if err != nil {
			return err
		}
		serviceContainers = append(serviceContainers, serviceContainer)
	}

	hostArtifactsFolder := fmt.Sprintf("%s/%s", paths.artifacts.host, vd.variant.ID)
	containerArtifactsFolder := fmt.Sprintf("%s/%s", paths.artifacts.container, vd.variant.ID)

	container, err := r.client.Run(&docker.RunOptions{
		Image: vd.imageTag,
		Links: containerLinks,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/var/run/docker.sock", paths.dockerSock.host),
			fmt.Sprintf("%s:/artifacts", hostArtifactsFolder),
		},
		Env: map[string]string{
			BazookaEnvSCM:           r.context.scm,
			BazookaEnvSCMUrl:        r.context.scmUrl,
			BazookaEnvSCMReference:  r.context.scmReference,
			BazookaEnvProjectID:     r.context.projectID,
			BazookaEnvJobID:         r.context.jobID,
			BazookaEnvJobParameters: r.context.jobParameters,
			"BZK_VARIANT":           strconv.Itoa(vd.variant.Number),
		},
		Detach:              true,
		LoggingDriver:       "syslog",
		LoggingDriverConfig: r.context.loggerConfig(vd.imageTag, vd.variant.ID),
	})
	if err != nil {
		return err
	}

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		if exitCode == 42 {
			return fmt.Errorf("Run failed\n Check Docker container logs, id is %s\n", container.ID())
		}
		success = false
	}
	if err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {
		return err
	}

	for _, serviceContainer := range serviceContainers {
		if err = serviceContainer.Remove(&docker.RemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		}); err != nil {
			return err
		}
	}

	if success {
		vd.variant.Status = commons.JOB_SUCCESS
	} else {
		vd.variant.Status = commons.JOB_FAILED
	}

	// Capture the artifacts list
	if err := filepath.Walk(containerArtifactsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		switch {
		case path == containerArtifactsFolder:
			return nil
		case info.IsDir():
			// nop
		default:
			relPath, err := filepath.Rel(containerArtifactsFolder, path)
			if err != nil {
				return err
			}
			vd.variant.Artifacts = append(vd.variant.Artifacts, relPath)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("Error while walking the artifacts: %v", err)
	}

	return nil
}

func safeDockerAlias(unsafeAlias string) string {
	re := regexp.MustCompile("(/|;|:|-|\\.)")
	return re.ReplaceAllString(unsafeAlias, "_")
}

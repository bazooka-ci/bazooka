package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	servicesFile := fmt.Sprintf("%s/%s/services", paths.work.container, vd.counter)

	servicesList, err := listServices(servicesFile)
	if err != nil {
		return err
	}

	serviceContainers := []*docker.Container{}
	containerLinks := []string{}
	for _, service := range servicesList {
		name := fmt.Sprintf("service-%s-%s-%d", r.context.projectID, r.context.jobID, vd.variant.Number)
		containerLinks = append(containerLinks, fmt.Sprintf("%s:%s", name, service))
		serviceContainer, err := r.client.Run(&docker.RunOptions{
			Name:   name,
			Image:  service,
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

func listServices(servicesFile string) ([]string, error) {
	file, err := os.Open(servicesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var services []string
	for scanner.Scan() {
		services = append(services, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

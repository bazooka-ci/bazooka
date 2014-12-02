package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	docker "github.com/bywan/go-dockercommand"
	commons "github.com/haklop/bazooka/commons"
	"github.com/haklop/bazooka/commons/mongo"
)

type Runner struct {
	BuildImages []BuiltImage
	Env         map[string]string
	Mongo       *mongo.MongoConnector
	client      *docker.Docker
}

func (r *Runner) Run(logger Logger) (bool, error) {
	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return false, err
	}
	r.client = client

	errChanRun := make(chan error)
	successChanRun := make(chan bool)
	remainingRuns := len(r.BuildImages)
	for _, buildImage := range r.BuildImages {
		go r.runContainer(logger, buildImage, r.Env, successChanRun, errChanRun)
	}

	success := true
	var lastError error
	for {
		select {
		case result := <-successChanRun:
			success = success && result
		case err := <-errChanRun:
			lastError = err
		}

		remainingRuns--
		if remainingRuns == 0 {
			break
		}
	}

	log.Printf("Dockerfiles builds finished\n")
	return success, lastError
}

func (r *Runner) runContainer(logger Logger, buildImage BuiltImage, env map[string]string, successChan chan bool, errChan chan error) {
	success := true
	servicesFile := fmt.Sprintf("%s/work/%d/services", BazookaInput, buildImage.VariantID)

	variant := &commons.Variant{
		Started:    time.Now(),
		BuildImage: buildImage.Image,
		Number:     buildImage.VariantID,
		JobID:      env[BazookaEnvJobID],
	}
	err := r.Mongo.AddVariant(variant)
	if err != nil {
		errChan <- err
		return
	}

	servicesList, err := listServices(servicesFile)
	if err != nil {
		mongoErr := r.Mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			errChan <- mongoErr
			return
		}
		errChan <- err
		return
	}

	serviceContainers := []*docker.Container{}
	containerLinks := []string{}
	for _, service := range servicesList {
		name := fmt.Sprintf("service-%s-%s-%d", env[BazookaEnvProjectID], env[BazookaEnvJobID], buildImage.VariantID)
		containerLinks = append(containerLinks, fmt.Sprintf("%s:%s", name, service))
		serviceContainer, err := r.client.Run(&docker.RunOptions{
			Name:   name,
			Image:  service,
			Detach: true,
		})
		if err != nil {
			mongoErr := r.Mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
			if mongoErr != nil {
				errChan <- mongoErr
				return
			}
			errChan <- err
			return
		}
		serviceContainers = append(serviceContainers, serviceContainer)
	}

	// TODO link containers
	container, err := r.client.Run(&docker.RunOptions{
		Image:  buildImage.Image,
		Links:  containerLinks,
		Detach: true,
	})
	if err != nil {
		mongoErr := r.Mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			errChan <- mongoErr
			return
		}
		errChan <- err
		return
	}

	container.Logs(buildImage.Image)
	logger(buildImage.Image, variant.ID, container)

	exitCode, err := container.Wait()
	if err != nil {
		mongoErr := r.Mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			errChan <- mongoErr
			return
		}
		errChan <- err
		return
	}
	if exitCode != 0 {
		if exitCode == 42 {
			mongoErr := r.Mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
			if mongoErr != nil {
				errChan <- mongoErr
				return
			}
			errChan <- fmt.Errorf("Run failed\n Check Docker container logs, id is %s\n", container.ID())
			return
		}
		success = false
	}
	err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil {
		mongoErr := r.Mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			errChan <- mongoErr
			return
		}
		errChan <- err
		return
	}

	for _, serviceContainer := range serviceContainers {
		err = serviceContainer.Remove(&docker.RemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil {
			mongoErr := r.Mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
			if mongoErr != nil {
				errChan <- mongoErr
				return
			}
			errChan <- err
			return
		}
	}
	var status commons.JobStatus
	if success {
		status = commons.JOB_SUCCESS
	} else {
		status = commons.JOB_FAILED
	}
	err = r.Mongo.FinishVariant(variant.ID, status, time.Now())
	if err != nil {
		errChan <- err
		return
	}
	successChan <- success
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

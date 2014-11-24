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
}

func (r *Runner) Run() (bool, error) {
	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return false, err
	}

	errChanRun := make(chan error)
	successChanRun := make(chan bool)
	remainingRuns := len(r.BuildImages)
	for _, buildImage := range r.BuildImages {
		go runContainer(client, buildImage, r.Env, r.Mongo, successChanRun, errChanRun)
	}

	success := false
	for {
		select {
		case result := <-successChanRun:
			success = success || result
			remainingRuns--
		case err := <-errChanRun:
			return false, err
		}

		if remainingRuns == 0 {
			break
		}
	}

	log.Printf("Dockerfiles builds finished\n")
	return success, nil
}

func runContainer(client *docker.Docker, buildImage BuiltImage, env map[string]string, mongo *mongo.MongoConnector, successChan chan bool, errChan chan error) {
	success := true
	servicesFile := fmt.Sprintf("%s/work/%d/services", BazookaInput, buildImage.VariantID)

	variant := &commons.Variant{
		Started:    time.Now(),
		BuildImage: buildImage.Image,
		Number:     buildImage.VariantID,
		JobID:      env[BazookaEnvJobID],
	}
	mongo.AddVariant(variant)

	servicesList, err := listServices(servicesFile)
	if err != nil {
		err2 := mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if err2 != nil {
			errChan <- err2
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
		serviceContainer, err := client.Run(&docker.RunOptions{
			Name:   name,
			Image:  service,
			Detach: true,
		})
		if err != nil {
			err2 := mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
			if err2 != nil {
				errChan <- err2
				return
			}
			errChan <- err
			return
		}
		serviceContainers = append(serviceContainers, serviceContainer)
	}

	// TODO link containers
	container, err := client.Run(&docker.RunOptions{
		Image:  buildImage.Image,
		Links:  containerLinks,
		Detach: true,
	})
	if err != nil {
		err2 := mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if err2 != nil {
			errChan <- err2
			return
		}
		errChan <- err
		return
	}

	container.Logs(buildImage.Image)

	exitCode, err := container.Wait()
	if err != nil {
		err2 := mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if err2 != nil {
			errChan <- err2
			return
		}
		errChan <- err
		return
	}
	if exitCode != 0 {
		if exitCode == 42 {
			err2 := mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
			if err2 != nil {
				errChan <- err2
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
		err2 := mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
		if err2 != nil {
			errChan <- err2
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
			err2 := mongo.FinishVariant(variant.ID, commons.JOB_ERRORED, time.Now())
			if err2 != nil {
				errChan <- err2
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
	err = mongo.FinishVariant(variant.ID, status, time.Now())
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

package main

import (
	"fmt"
	"log"

	docker "github.com/bywan/go-dockercommand"
)

type Runner struct {
	BuildImages []string
}

func (r *Runner) Run() error {
	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return err
	}

	errChanRun := make(chan error)
	successChanRun := make(chan bool)
	remainingRuns := len(r.BuildImages)
	for _, buildImage := range r.BuildImages {
		go runContainer(client, buildImage, successChanRun, errChanRun)
	}

	for {
		select {
		case _ = <-successChanRun:
			remainingRuns--
		case err := <-errChanRun:
			return err
		}

		if remainingRuns == 0 {
			break
		}
	}

	log.Printf("Dockerfiles builds finished\n")
	return nil
}

func runContainer(client *docker.Docker, buildImage string, successChan chan bool, errChan chan error) {
	container, err := client.Run(&docker.RunOptions{
		Image:  buildImage,
		Detach: true,
	})
	if err != nil {
		errChan <- err
		return
	}

	container.Logs(buildImage)

	exitCode, err := container.Wait()
	if err != nil {
		errChan <- err
		return
	}
	if exitCode != 0 {
		errChan <- fmt.Errorf("Run failed\n Check Docker container logs, id is %s\n", container.ID())
		return
	}
	successChan <- true
}

package main

import (
	"fmt"
	"log"
	"os"

	docker "github.com/bywan/go-dockercommand"
)

const (
	dockerEndpoint = "unix:///docker.sock"
)

type LanguageParser struct {
	Image string
}

func (p *LanguageParser) Parse() error {

	log.Printf("Lauching language parser %s\n", p.Image)

	client, err := docker.NewDocker(dockerEndpoint)
	if err != nil {
		return err
	}
	bazookaHome := os.Getenv("BZK_HOME")
	container, err := client.Run(&docker.RunOptions{
		Image: p.Image,
		VolumeBinds: []string{
			fmt.Sprintf("%s/source/:/bazooka", bazookaHome),
			fmt.Sprintf("%s/work/:/bazooka-output", bazookaHome)},
	})
	if err != nil {
		return err
	}

	details, err := container.Inspect()
	if err != nil {
		return err
	}
	if details.State.ExitCode != 0 {
		return fmt.Errorf("Error during execution of Language Parser container %s/parser\n Check Docker container logs, id is %s\n", p.Image, container.ID())
	}

	return container.Remove()
}

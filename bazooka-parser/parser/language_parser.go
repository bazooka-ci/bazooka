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
	Options *LanguageParseOptions
}

type LanguageParseOptions struct {
	InputFolder string
	Image       string
}

func (p *LanguageParser) Parse() error {

	log.Printf("Lauching language parser %s in %s\n", p.Options.Image, p.Options.InputFolder)

	client, err := docker.NewDocker(dockerEndpoint)
	if err != nil {
		return err
	}
	bazookaHome := os.Getenv("BZK_HOME")
	containerID, err := client.Run(&docker.RunOptions{
		Image: p.Options.Image,
		VolumeBinds: []string{
			fmt.Sprintf("%s/source/:/bazooka", bazookaHome),
			fmt.Sprintf("%s/work/:/bazooka-output", bazookaHome)},
	})
	if err != nil {
		return err
	}

	details, err := client.Inspect(containerID)
	if err != nil {
		return err
	}
	if details.State.ExitCode != 0 {
		return fmt.Errorf("Error during execution of Language Parser container %s/parser\n Check Docker container logs, id is %s\n", p.Options.Image, containerID)
	}

	return client.Rm(&docker.RmOptions{
		Container: []string{containerID},
	})
}

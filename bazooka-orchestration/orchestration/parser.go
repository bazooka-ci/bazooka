package main

import (
	"fmt"
	"log"

	docker "github.com/bywan/go-dockercommand"
)

const (
	BazookaParseImage = "bazooka/parser"
)

type Parser struct {
	Options *ParseOptions
}

type ParseOptions struct {
	InputFolder    string
	OutputFolder   string
	DockerSock     string
	HostBaseFolder string
	Env            map[string]string
}

func (p *Parser) Parse() error {

	log.Printf("Parsing Configuration from checked-out source %s\n", p.Options.InputFolder)
	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return err
	}
	containerID, err := client.Run(&docker.RunOptions{
		Image: BazookaParseImage,
		Env:   p.Options.Env,
		VolumeBinds: []string{fmt.Sprintf("%s:/bazooka", p.Options.InputFolder), fmt.Sprintf("%s:/bazooka-output", p.Options.OutputFolder),
			fmt.Sprintf("%s:/docker.sock", p.Options.DockerSock)},
	})
	if err != nil {
		return err
	}

	details, err := client.Inspect(containerID)
	if err != nil {
		return err
	}
	if details.State.ExitCode != 0 {
		return fmt.Errorf("Error during execution of Parser container %s/parser\n Check Docker container logs, id is %s\n", BazookaParseImage, containerID)
	}

	err = client.Rm(&docker.RmOptions{
		Container: []string{containerID},
	})
	if err != nil {
		return err
	}
	log.Printf("Configuration parsed and Dockerfiles generated in %s\n", p.Options.OutputFolder)
	return nil
}

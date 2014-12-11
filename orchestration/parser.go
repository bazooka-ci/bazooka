package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
	"github.com/haklop/bazooka/commons/mongo"
)

const (
	BazookaParseImage = "bazooka/parser"
)

type Parser struct {
	MongoConnector *mongo.MongoConnector
	Options        *ParseOptions
}

type ParseOptions struct {
	InputFolder    string
	OutputFolder   string
	DockerSock     string
	HostBaseFolder string
	MetaFolder     string
	Env            map[string]string
}

func (p *Parser) Parse(logger Logger) error {

	log.Info("Running Parsing Image %s on checked-out source\n", BazookaParseImage)
	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return err
	}

	image, err := p.resolveImage()
	if err != nil {
		return err
	}

	log.Info("Using image '%s'\n", image)

	container, err := client.Run(&docker.RunOptions{
		Image: image,
		Env:   p.Options.Env,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", p.Options.InputFolder),
			fmt.Sprintf("%s:/meta", p.Options.MetaFolder),
			fmt.Sprintf("%s:/bazooka-output", p.Options.OutputFolder),
			fmt.Sprintf("%s:/docker.sock", p.Options.DockerSock)},
		Detach: true,
	})
	if err != nil {
		return err
	}

	container.Logs(BazookaParseImage)
	logger(BazookaParseImage, "", container)

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("Error during execution of Parser container %s/parser\n Check Docker container logs, id is %s\n", BazookaParseImage, container.ID())
	}

	err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil {
		return err
	}
	log.Info("Parsing Image ran sucessfully, Dockerfiles generated in %s\n", p.Options.OutputFolder)
	return nil
}

func (f *Parser) resolveImage() (string, error) {
	image, err := f.MongoConnector.GetImage("parser")
	if err != nil {
		return "", fmt.Errorf("Unable to find Bazooka Docker Image for parser\n")
	}
	return image, nil
}

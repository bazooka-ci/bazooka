package main

import (
	"fmt"

	docker "github.com/bywan/go-dockercommand"
	lib "github.com/haklop/bazooka/commons"
	l "github.com/haklop/bazooka/commons/logger"
	"github.com/haklop/bazooka/commons/mongo"
)

type SCMFetcher struct {
	Options        *FetchOptions
	MongoConnector *mongo.MongoConnector
}

type FetchOptions struct {
	Scm         string
	URL         string
	Reference   string
	LocalFolder string
	KeyFile     string
	MetaFolder  string
	JobID       string
	Env         map[string]string
}

func (f *SCMFetcher) Fetch(logger Logger) error {

	l.Info.Printf("Fetching SCM From Source Repository at %s\n", f.Options.URL)

	image, err := f.resolveImage()
	if err != nil {
		return err
	}
	l.Info.Printf("Using image '%s'\n", image)

	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return err
	}
	container, err := client.Run(&docker.RunOptions{
		Image: image,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", f.Options.LocalFolder),
			fmt.Sprintf("%s:/bazooka-key", f.Options.KeyFile),
			fmt.Sprintf("%s:/meta", f.Options.MetaFolder),
		},
		Env:    f.Options.Env,
		Detach: true,
	})
	if err != nil {
		return err
	}

	container.LogsWith(image, l.Docker)
	logger(image, "", container)

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("Error during execution of SCM container %s\n Check Docker container logs, id is %s\n", image, container.ID())
	}

	l.Info.Printf("SCM Source Repo Fetched in %s\n", f.Options.LocalFolder)
	err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	scmMetadata := &lib.SCMMetadata{}
	localMetaFolder := fmt.Sprintf(MetaFolderPattern, BazookaInput)
	scmMetadataFile := fmt.Sprintf("%s/scm", localMetaFolder)
	err = lib.Parse(scmMetadataFile, scmMetadata)
	if err != nil {
		return err
	}

	err = f.MongoConnector.AddJobSCMMetadata(f.Options.JobID, scmMetadata)
	if err != nil {
		return err
	}
	return err
}

func (f *SCMFetcher) resolveImage() (string, error) {
	image, err := f.MongoConnector.GetImage(fmt.Sprintf("scm/fetch/%s", f.Options.Scm))
	if err != nil {
		return "", fmt.Errorf("Unable to find Bazooka Docker Image for SCM %s\n", f.Options.Scm)
	}
	return image, nil
}

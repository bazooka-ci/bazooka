package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/commons/mongo"
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

	log.WithFields(log.Fields{
		"source": f.Options.URL,
	}).Info("Fetching SCM From Source Repository")

	image, err := f.resolveImage()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"image": image,
	}).Info("Starting SCM Fetch")

	client, err := docker.NewDocker(paths.container.dockerEndpoint)
	if err != nil {
		return err
	}

	volumes := []string{
		fmt.Sprintf("%s:/bazooka", f.Options.LocalFolder),
		fmt.Sprintf("%s:/meta", f.Options.MetaFolder),
	}
	if len(f.Options.KeyFile) > 0 {
		volumes = append(volumes, fmt.Sprintf("%s:/bazooka-key", f.Options.KeyFile))
	}

	container, err := client.Run(&docker.RunOptions{
		Image:       image,
		VolumeBinds: volumes,
		Env:         f.Options.Env,
		Detach:      true,
	})
	if err != nil {
		return err
	}

	container.Logs(image)
	logger(image, "", container)

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("Error during execution of SCM container %s\n Check Docker container logs, id is %s\n", image, container.ID())
	}

	log.WithFields(log.Fields{
		"checkout_folder": f.Options.LocalFolder,
	}).Info("SCM Source Repo Fetched")

	err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	scmMetadata := &lib.SCMMetadata{}
	
	scmMetadataFile := fmt.Sprintf("%s/scm", paths.container.meta)
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

package main

import (
	"fmt"
	"log"

	docker "github.com/bywan/go-dockercommand"
	lib "github.com/haklop/bazooka/commons"
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

	log.Printf("Fetching SCM From Source Repo %s\n", f.Options.URL)

	image, err := resolveSCMImage(f.Options.Scm)
	if err != nil {
		return err
	}

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

	container.Logs(image)
	logger(image, "", container)

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("Error during execution of SCM container %s\n Check Docker container logs, id is %s\n", image, container.ID())
	}

	log.Printf("SCM Source Repo Fetched in %s\n", f.Options.LocalFolder)
	err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	scmMetadata := &lib.SCMMetadata{}
	localMetaFolder := fmt.Sprintf(MetaFolderPattern, BazookaInput)
	scmMetadataFile := fmt.Sprintf("%s/scm", localMetaFolder)
	log.Printf("Parsing SCM Metadata in %s\n", scmMetadataFile)
	err = lib.Parse(scmMetadataFile, scmMetadata)
	if err != nil {
		return err
	}
	log.Printf("Metadata Parsed is %+v\n", scmMetadata)

	f.MongoConnector.AddJobSCMMetadata(f.Options.JobID, scmMetadata)
	return err
}

func resolveSCMImage(scm string) (string, error) {
	//TODO extract this from db
	scmMap := map[string]string{
		"git": "bazooka/scm-git",
	}
	if val, ok := scmMap[scm]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Unable to find Bazooka Docker Image for SCM %s\n", scm)
}

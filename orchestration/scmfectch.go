package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	docker "github.com/bywan/go-dockercommand"
)

type SCMFetcher struct {
	context *context
	update  bool
}

func (f *SCMFetcher) Fetch() error {
	log.WithFields(log.Fields{
		"source": f.context.scmUrl,
	}).Info("Fetching SCM From Source Repository")

	paths := f.context.paths

	image, err := f.resolveImage()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"image": image,
	}).Info("Starting SCM Fetch")

	client, err := docker.NewDocker(paths.dockerEndpoint.container)
	if err != nil {
		return err
	}

	env := map[string]string{
		BazookaEnvSCMUrl:       f.context.scmUrl,
		BazookaEnvSCMReference: f.context.scmReference,
		BazookaEnvProjectID:    f.context.projectID,
		BazookaEnvJobID:        f.context.jobID,
	}

	if f.update {
		env["UPDATE"] = "1"
	}

	volumes := []string{
		fmt.Sprintf("%s:/bazooka", paths.source.host),
		fmt.Sprintf("%s:/meta", paths.meta.host),
	}
	scmKeyFile := paths.scmKey.host
	if len(scmKeyFile) > 0 {
		volumes = append(volumes, fmt.Sprintf("%s:/bazooka-key", scmKeyFile))
	}

	container, err := client.Run(&docker.RunOptions{
		Image:               image,
		VolumeBinds:         volumes,
		Env:                 env,
		Detach:              true,
		NetworkMode:         f.context.network,
		LoggingDriver:       "syslog",
		LoggingDriverConfig: f.context.loggerConfig(image, ""),
	})
	if err != nil {
		return fmt.Errorf("Failed to run the scm container: %v", err)
	}

	defer lib.RemoveContainer(container)

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("Error during execution of SCM container %s, exit code %d\n Check Docker container logs, id is %s\n",
			image, exitCode, container.ID())
	}

	log.WithFields(log.Fields{
		"checkout_folder": paths.source.host,
	}).Info("SCM Source Repo Fetched")

	scmMetadata := &lib.SCMMetadata{}

	scmMetadataFile := fmt.Sprintf("%s/scm", paths.meta.container)
	err = lib.Parse(scmMetadataFile, scmMetadata)
	if err != nil {
		return err
	}

	err = f.context.client.Internal.AddJobSCMMetadata(f.context.jobID, scmMetadata)
	if err != nil {
		return err
	}
	return err
}

func (f *SCMFetcher) resolveImage() (string, error) {
	image, err := f.context.client.Image.Get(fmt.Sprintf("scm/fetch/%s", f.context.scm))
	if err != nil {
		return "", fmt.Errorf("Unable to find Bazooka Docker Image for SCM %s\n, error is %v", f.context.scm, err)
	}
	return image.Image, nil
}

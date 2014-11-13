package main

import (
	"fmt"
	"log"

	docker "github.com/bywan/go-dockercommand"
)

type SCMFetcher struct {
	Options *FetchOptions
}

type FetchOptions struct {
	Scm         string
	URL         string
	Reference   string
	LocalFolder string
	KeyFile     string
	Env         map[string]string
}

func (f *SCMFetcher) Fetch() error {

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
		Image:       image,
		VolumeBinds: []string{fmt.Sprintf("%s:/bazooka", f.Options.LocalFolder), fmt.Sprintf("%s:/bazooka-key", f.Options.KeyFile)},
		Env:         f.Options.Env,
		Detach:      true,
	})
	if err != nil {
		return err
	}

	container.Logs(image)

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("Error during execution of SCM container %s\n Check Docker container logs, id is %s\n", image, container.ID())
	}

	log.Printf("SCM Source Repo Fetched in %s\n", f.Options.LocalFolder)
	return container.Remove()
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

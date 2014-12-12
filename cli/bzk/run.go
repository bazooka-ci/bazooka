package main

import (
	"fmt"
	"log"
	"strings"

	docker "github.com/bywan/go-dockercommand"
	"github.com/codegangsta/cli"
	dockerclient "github.com/fsouza/go-dockerclient"
)

var allContainers []dockerclient.APIContainers
var registry string

func run(c *cli.Context) {
	bzkHome := c.String("home")
	if len(bzkHome) == 0 {
		log.Fatal("$BZK_HOME environment variable is needed (or use --bzk-home option)")
	}
	scmKey := c.String("scm-key")
	if len(scmKey) == 0 {
		log.Fatal("$BZK_SCM_KEYFILE environment variable is needed (or use --scm-key option)")
	}
	forceRestart := c.Bool("restart")
	forceUpdate := c.Bool("update")
	registry = c.String("registry")

	client, err := docker.NewDocker("")
	if err != nil {
		log.Fatal(err)
	}

	allContainers, err = client.Ps(&docker.PsOptions{
		All: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if forceUpdate {
		log.Printf("Pulling Bazooka images to check for new versions\n")
		mandatoryImages := []string{"server", "web", "orchestration", "parser"}
		optionalImages := []string{"parser-java", "parser-golang", "scm-git",
			"runner-java", "runner-java:oraclejdk8", "runner-java:oraclejdk7", "runner-java:oraclejdk6", "runner-java:openjdk8", "runner-java:openjdk7", "runner-java:openjdk6",
			"runner-golang", "runner-golang:1.2.2", "runner-golang:1.3", "runner-golang:1.3.1", "runner-golang:1.3.2", "runner-golang:1.3.3", "runner-golang:1.4"}
		for _, image := range mandatoryImages {
			err = client.Pull(&docker.PullOptions{Image: getImageLocation(fmt.Sprintf("bazooka/%s", image))})
			if err != nil {
				log.Fatal(fmt.Errorf("Unable to pull required image for Bazooka, reason is: %v\n", err))
			}
		}
		for _, image := range optionalImages {
			err = client.Pull(&docker.PullOptions{Image: getImageLocation(fmt.Sprintf("bazooka/%s", image))})
			if err != nil {
				log.Printf("Unable to pull image for Bazooka, as it is an optional one, let's move on. Reason is: %v\n", err)
			}
		}
	}

	mongoRestarted, err := ensureContainerIsRestarted(client, &docker.RunOptions{
		Name: "bzk_mongodb",
		// Using the official mongo image from dockerhub, this may need a change later
		Image:  "mongo",
		Detach: true,
	}, false)

	serverRestarted, err := ensureContainerIsRestarted(client, &docker.RunOptions{
		Name:   "bzk_server",
		Image:  getImageLocation("bazooka/server"),
		Detach: true,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", c.String("home")),
			fmt.Sprintf("%s:/var/run/docker.sock", c.String("docker-sock")),
		},
		Links: []string{"bzk_mongodb:mongo"},
		Env: map[string]string{
			"BZK_SCM_KEYFILE": c.String("scm-key"),
			"BZK_HOME":        c.String("home"),
			"BZK_DOCKERSOCK":  c.String("docker-sock"),
		},
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"3000/tcp": []dockerclient.PortBinding{
				dockerclient.PortBinding{HostPort: "3000"},
			},
		},
	}, mongoRestarted || forceRestart || forceUpdate)

	_, err = ensureContainerIsRestarted(client, &docker.RunOptions{
		Name:   "bzk_web",
		Image:  getImageLocation("bazooka/web"),
		Detach: true,
		Links:  []string{"bzk_server:server"},
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"80/tcp": []dockerclient.PortBinding{
				dockerclient.PortBinding{HostPort: "8000"},
			},
		},
	}, serverRestarted || forceRestart || forceUpdate)

	if err != nil {
		log.Fatal(err)
	}

}

func ensureContainerIsRestarted(client *docker.Docker, options *docker.RunOptions, needRestart bool) (bool, error) {
	container, err := getContainer(allContainers, options.Name)
	if err != nil {
		log.Printf("Container %s not found, Starting it\n", options.Name)
		_, err := client.Run(options)
		return true, err
	}
	if needRestart {
		log.Printf("Restarting Container %s\n", options.Name)
		err = client.Rm(&docker.RmOptions{
			Container: []string{container.ID},
			Force:     true,
		})
		if err != nil {
			return false, err
		}
		_, err := client.Run(options)
		return true, err
	}
	if strings.HasPrefix(container.Status, "Up") {
		log.Printf("Container %s already Up & Running, keeping on\n", options.Name)
		return false, nil
	}
	log.Printf("Container %s is not `Up`, starting it\n", options.Name)
	return true, client.Start(&docker.StartOptions{
		ID: container.ID,
	})

}

func getContainer(containers []dockerclient.APIContainers, name string) (dockerclient.APIContainers, error) {
	for _, container := range containers {
		if contains(container.Names, name) || contains(container.Names, "/"+name) {
			return container, nil
		}
	}
	return dockerclient.APIContainers{}, fmt.Errorf("Container not found")
}

func getImageLocation(image string) string {
	if len(registry) > 0 {
		return fmt.Sprintf("%s/%s", registry, image)
	}
	return image
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

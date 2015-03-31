package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jawher/mow.cli"

	docker "github.com/bywan/go-dockercommand"

	dockerclient "github.com/fsouza/go-dockerclient"
)

var allContainers []dockerclient.APIContainers

func run(cmd *cli.Cmd) {
	cmd.Spec = "[--home|--scm-key|--registry]... [--restart|--update]"

	bzkHome := cmd.String(cli.StringOpt{
		Name:   "home",
		Desc:   "Bazooka's work directory",
		EnvVar: "BZK_HOME",
	})
	scmKey := cmd.String(cli.StringOpt{
		Name:   "scm-key",
		Desc:   "Location of the private SSH Key Bazooka will use for SCM Fetch",
		EnvVar: "BZK_SCM_KEYFILE",
	})
	registry := cmd.String(cli.StringOpt{
		Name:   "registry",
		EnvVar: "BZK_REGISTRY",
	})
	dockerSock := cmd.String(cli.StringOpt{
		Name:   "docker-sock",
		Desc:   "Location of the Docker unix socket, usually /var/run/docker.sock",
		EnvVar: "BZK_DOCKERSOCK",
	})

	forceRestart := cmd.Bool(cli.BoolOpt{
		Name: "r restart",
		Desc: "Restart Bazooka if already running",
	})
	forceUpdate := cmd.Bool(cli.BoolOpt{
		Name: "u update",
		Desc: "Update Bazooka to the latest version by pulling new images from the registry",
	})

	cmd.Action = func() {
		config, err := loadConfig()
		if err != nil {
			log.Fatal(fmt.Errorf("Unable to load Bazooka config, reason is: %v\n", err))
		}
		if len(*bzkHome) == 0 {
			if len(config.Home) == 0 {
				*bzkHome = interactiveInput("Bazooka Home Folder")
				config.Home = *bzkHome
			} else {
				*bzkHome = config.Home
			}
		}

		if len(*dockerSock) == 0 {
			if len(config.DockerSock) == 0 {
				*dockerSock = interactiveInput("Docker Socket path")
				config.DockerSock = *dockerSock
			} else {
				*dockerSock = config.DockerSock
			}
		}

		if len(*scmKey) == 0 {
			if len(config.SCMKey) == 0 {
				*scmKey = interactiveInput("Bazooka Default SCM private key")
				config.SCMKey = *scmKey
			} else {
				*scmKey = config.SCMKey
			}
		}

		err = saveConfig(config)
		if err != nil {
			log.Fatal(fmt.Errorf("Unable to save Bazooka config, reason is: %v\n", err))
		}

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

		if *forceUpdate {
			log.Printf("Pulling Bazooka images to check for new versions\n")
			mandatoryImages := []string{"server", "web", "orchestration", "parser"}
			optionalImages := []string{"parser-java", "parser-golang", "scm-git",
				"runner-java", "runner-java:oraclejdk8", "runner-java:oraclejdk7", "runner-java:oraclejdk6", "runner-java:openjdk8", "runner-java:openjdk7", "runner-java:openjdk6",
				"runner-golang", "runner-golang:1.2.2", "runner-golang:1.3", "runner-golang:1.3.1", "runner-golang:1.3.2", "runner-golang:1.3.3", "runner-golang:1.4"}
			for _, image := range mandatoryImages {
				err = client.Pull(&docker.PullOptions{Image: getImageLocation(*registry, fmt.Sprintf("bazooka/%s", image))})
				if err != nil {
					log.Fatal(fmt.Errorf("Unable to pull required image for Bazooka, reason is: %v\n", err))
				}
			}
			for _, image := range optionalImages {
				err = client.Pull(&docker.PullOptions{Image: getImageLocation(*registry, fmt.Sprintf("bazooka/%s", image))})
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
			Image:  getImageLocation(*registry, "bazooka/server"),
			Detach: true,
			VolumeBinds: []string{
				fmt.Sprintf("%s:/bazooka", *bzkHome),
				fmt.Sprintf("%s:/var/run/docker.sock", *dockerSock),
			},
			Links: []string{"bzk_mongodb:mongo"},
			Env:   getServerEnv(*bzkHome, *dockerSock, *scmKey),
			PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
				"3000/tcp": {{HostPort: "3000"}},
			},
		}, mongoRestarted || *forceRestart || *forceUpdate)

		_, err = ensureContainerIsRestarted(client, &docker.RunOptions{
			Name:   "bzk_web",
			Image:  getImageLocation(*registry, "bazooka/web"),
			Detach: true,
			Links:  []string{"bzk_server:server"},
			PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
				"80/tcp": {{HostPort: "8000"}},
			},
		}, serverRestarted || *forceRestart || *forceUpdate)

		if err != nil {
			log.Fatal(err)
		}
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

func getServerEnv(home, dockerSock, scmKey string) map[string]string {
	envMap := map[string]string{
		"BZK_HOME":       home,
		"BZK_DOCKERSOCK": dockerSock,
	}
	if len(scmKey) > 0 {
		envMap["BZK_SCM_KEYFILE"] = scmKey
	}
	return envMap
}

func getContainer(containers []dockerclient.APIContainers, name string) (dockerclient.APIContainers, error) {
	for _, container := range containers {
		if contains(container.Names, name) || contains(container.Names, "/"+name) {
			return container, nil
		}
	}
	return dockerclient.APIContainers{}, fmt.Errorf("Container not found")
}

func getImageLocation(registry, image string) string {
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

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

func startService(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"

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
	tag := cmd.String(cli.StringOpt{
		Name: "tag",
		Desc: "The bazooka version to run",
	})

	cmd.Action = func() {
		config, err := getConfigWithParams(*bzkHome, *dockerSock, *registry, *scmKey)
		if err != nil {
			log.Fatal(err)
		}

		if len(*tag) == 0 {
			*tag = "latest"
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

		err = ensureContainerIsStarted(client, getMongoRunOptions())
		if err != nil {
			log.Fatal(err)
		}

		err = ensureContainerIsStarted(client, getServerRunOptions(config.Registry, config.Home, config.DockerSock, config.SCMKey, *tag))
		if err != nil {
			log.Fatal(err)
		}

		err = ensureContainerIsStarted(client, getWebRunOptions(config.Registry, *tag))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func restartService(cmd *cli.Cmd) {
	cmd.Action = doRestartService
}

func doRestartService() {
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

	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = ensureContainerIsStarted(client, getMongoRunOptions())
	if err != nil {
		log.Fatal(err)
	}

	tag, err := getTagFromCurrentImages(client)

	err = restartContainer(client, getServerRunOptions(config.Registry, config.Home, config.DockerSock, config.SCMKey, tag))
	if err != nil {
		log.Fatal(err)
	}

	err = restartContainer(client, getWebRunOptions(config.Registry, tag))
	if err != nil {
		log.Fatal(err)
	}
}

func getTagFromCurrentImages(client *docker.Docker) (string, error) {
	container, err := getContainer(allContainers, "bzk_server")
	if err != nil {
		// Container not found, using latest by default
		return "latest", nil
	}
	split := strings.Split(container.Image, ":")
	return split[len(split)-1], nil
}

func upgradeService(cmd *cli.Cmd) {
	cmd.Action = func() {
		fmt.Printf("Pulling Bazooka images to check for new versions. This may take some time...\n")

		client, err := docker.NewDocker("")
		if err != nil {
			log.Fatal(err)
		}

		bzkImages, err := client.Images(&docker.ImagesOptions{})
		if err != nil {
			log.Fatal(err)
		}

		var hasError bool
		for _, image := range bzkImages {
			for _, tag := range image.RepoTags {
				if strings.HasPrefix(tag, "bazooka/") {
					err = client.Pull(&docker.PullOptions{Image: tag})
					if err != nil {
						fmt.Printf("Unable to pull image %s, reason is: %v\n", tag, err)
						hasError = true
					}
					fmt.Printf("Newest image %s upgraded\n", tag)
					break
				}
			}
		}
		if hasError {
			log.Fatalf("Error while pulling some images, see above errors\n")
		}

		fmt.Printf("Bazooka Images have been upgraded, let's restart bazooka\n")
		doRestartService()
	}
}

func stopService(cmd *cli.Cmd) {
	cmd.Action = func() {
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

		err = stopContainer(client, "bzk_mongodb")
		if err != nil {
			log.Fatal(err)
		}

		err = stopContainer(client, "bzk_server")
		if err != nil {
			log.Fatal(err)
		}

		err = stopContainer(client, "bzk_web")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func statusService(cmd *cli.Cmd) {
	cmd.Action = func() {
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

		mongoUp := getContainerStatus("bzk_mongodb")
		serverUp := getContainerStatus("bzk_server")
		webUp := getContainerStatus("bzk_web")

		if mongoUp && serverUp && webUp {
			fmt.Printf("Bazooka service is Up\n")
		} else {
			fmt.Printf("Bazooka service is Down\n")
		}

	}
}

func getContainerStatus(name string) bool {
	container, err := getContainer(allContainers, name)
	if err != nil {
		fmt.Printf("Container %s not started\n", name)
		return false
	}
	if strings.HasPrefix(container.Status, "Up") {
		split := strings.Split(container.Image, ":")
		fmt.Printf("Container %s running, version \"%s\"\n", name, split[len(split)-1])
		return true
	}
	fmt.Printf("Container %s stopped\n", name)
	return false

}

func getConfigWithParams(bzkHome, dockerSock, registry, scmKey string) (*Config, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to load Bazooka config, reason is: %v\n", err)
	}
	if len(bzkHome) == 0 {
		if len(config.Home) == 0 {
			config.Home = interactiveInput("Bazooka Home Folder")
		}
	} else {
		config.Home = bzkHome
	}
	if len(dockerSock) == 0 {
		if len(config.DockerSock) == 0 {
			config.DockerSock = interactiveInput("Docker Socket path")
		}
	} else {
		config.DockerSock = dockerSock
	}

	if len(scmKey) == 0 {
		if len(config.SCMKey) == 0 {
			config.SCMKey = interactiveInput("Bazooka Default SCM private key")
		}
	} else {
		config.SCMKey = scmKey
	}

	if len(registry) != 0 {
		config.Registry = registry
	}

	err = saveConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Unable to save Bazooka config, reason is: %v\n", err)
	}

	return config, nil
}

func getConfig() (*Config, error) {
	return getConfigWithParams("", "", "", "")
}

func restartContainer(client *docker.Docker, options *docker.RunOptions) error {
	container, err := getContainer(allContainers, options.Name)
	if err != nil {
		fmt.Printf("Container %s not found, Starting it\n", options.Name)
		_, err := client.Run(options)
		return err
	}
	fmt.Printf("Restarting Container %s\n", options.Name)
	err = client.Rm(&docker.RmOptions{
		Container: []string{container.ID},
		Force:     true,
	})
	if err != nil {
		return err
	}
	_, err = client.Run(options)
	return err
}

func stopContainer(client *docker.Docker, name string) error {
	container, err := getContainer(allContainers, name)
	if err != nil {
		fmt.Printf("Container %s not found, doing nothing\n", name)
		return nil
	}
	fmt.Printf("Stopping Container %s\n", name)
	err = client.Stop(&docker.StopOptions{
		ID: container.ID,
	})
	if err != nil {
		return fmt.Errorf("Error stopping container %s, reason is %v\n", name, err)
	}
	return nil
}

func ensureContainerIsStarted(client *docker.Docker, options *docker.RunOptions) error {
	container, err := getContainer(allContainers, options.Name)
	if err != nil {
		fmt.Printf("Container %s not found, Starting it\n", options.Name)
		_, err := client.Run(options)
		return err
	}
	if container.Image == options.Image {
		if strings.HasPrefix(container.Status, "Up") {
			fmt.Printf("Container %s already Up & Running, keeping on\n", options.Name)
			return nil
		}
		fmt.Printf("Container %s is not `Up`, starting it\n", options.Name)
		return client.Start(&docker.StartOptions{
			ID: container.ID,
		})
	}
	fmt.Printf("Container %s found, but not in the right version, recreating it\n", options.Name)
	return restartContainer(client, options)
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

func getImageLocation(registry, image, tag string) string {
	location := image
	if len(registry) > 0 {
		location = fmt.Sprintf("%s/%s", registry, image)
	}
	if len(tag) > 0 {
		location = fmt.Sprintf("%s:%s", location, tag)
	}
	return location
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func getMongoRunOptions() *docker.RunOptions {
	return &docker.RunOptions{
		Name: "bzk_mongodb",
		// Using the official mongo image from dockerhub, this may need a change later
		Image:  "mongo:3.0.2",
		Detach: true,
	}
}

func getServerRunOptions(registry, bzkHome, dockerSock, scmKey, tag string) *docker.RunOptions {
	return &docker.RunOptions{
		Name:   "bzk_server",
		Image:  getImageLocation(registry, "bazooka/server", tag),
		Detach: true,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", bzkHome),
			fmt.Sprintf("%s:/var/run/docker.sock", dockerSock),
		},
		Links: []string{"bzk_mongodb:mongo"},
		Env:   getServerEnv(bzkHome, dockerSock, scmKey),
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"3000/tcp": {{HostPort: "3000"}},
		},
	}
}

func getWebRunOptions(registry, tag string) *docker.RunOptions {
	return &docker.RunOptions{
		Name:   "bzk_web",
		Image:  getImageLocation(registry, "bazooka/web", tag),
		Detach: true,
		Links:  []string{"bzk_server:server"},
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"80/tcp": {{HostPort: "8000"}},
		},
	}
}

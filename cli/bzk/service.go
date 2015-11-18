package main

import (
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/jawher/mow.cli"

	docker "github.com/bywan/go-dockercommand"

	dockerclient "github.com/fsouza/go-dockerclient"
)

const (
	bzkNetwork         = "bzk_net"
	bzkContainerMongo  = "bzk_db"
	bzkContainerServer = "bzk_server"
	bzkContainerWeb    = "bzk_web"
)

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
	syslogURL := cmd.String(cli.StringOpt{
		Name:   "syslog-url",
		Desc:   "URL of the bazooka server syslog endpoint",
		EnvVar: "BZK_SYSLOG_URL",
	})
	dbURL := cmd.String(cli.StringOpt{
		Name:   "db-url",
		Desc:   "URL of a MongoDB server",
		EnvVar: "BZK_DB_URL",
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

	cmd.Action = doStartService(tag, bzkHome, dockerSock, registry, scmKey, syslogURL, dbURL)
}

func doStartService(version, bzkHome, dockerSock, registry, scmKey, syslogURL, dbURL *string) func() {
	return func() {
		ensureNetworkExists()

		config, err := getConfigWithParams(*version, *bzkHome, *dockerSock, *registry, *scmKey, *syslogURL, *dbURL)
		if err != nil {
			log.Fatal(err)
		}

		client, err := docker.NewDocker("")
		if err != nil {
			log.Fatal(err)
		}

		allContainers, err := client.Ps(&docker.PsOptions{
			All: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		if config.DbURL == "" {
			err = startContainer(client, getMongoRunOptions(), allContainers)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = startContainer(client, getServerRunOptions(config), allContainers)
		if err != nil {
			log.Fatal(err)
		}

		err = startContainer(client, getWebRunOptions(config), allContainers)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ensureNetworkExists() {
	client, err := docker.NewDocker("")
	if err != nil {
		log.Fatal(err)
	}

	networks, err := client.Networks()
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range networks {
		if n.Name == bzkNetwork {
			return
		}
	}

	if _, err := client.CreateNetwork(dockerclient.CreateNetworkOptions{
		Name:   bzkNetwork,
		Driver: "bridge",
	}); err != nil {
		log.Fatalf("Error while creating bridge network %s: %v", bzkNetwork, err)
	}
}

func stopService(cmd *cli.Cmd) {
	cmd.Action = doStopService()
}

func doStopService() func() {
	return func() {
		client, err := docker.NewDocker("")
		if err != nil {
			log.Fatal(err)
		}

		config, err := getConfig()
		if err != nil {
			log.Fatal(err)
		}

		allContainers, err := client.Ps(&docker.PsOptions{
			All: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		if config.DbURL == "" {
			err := stopContainer(client, bzkContainerMongo, allContainers)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = stopContainer(client, bzkContainerServer, allContainers)
		if err != nil {
			log.Fatal(err)
		}

		err = stopContainer(client, bzkContainerWeb, allContainers)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func restartService(cmd *cli.Cmd) {
	cmd.Spec = "[-x [-swd]]"

	recreate := cmd.Bool(cli.BoolOpt{
		Name:      "x recreate",
		Desc:      "Recreate existing Bazooka containers",
		Value:     false,
		HideValue: true,
	})
	recreateServer := cmd.Bool(cli.BoolOpt{
		Name:      "s server",
		Desc:      "Recreate existing server container",
		Value:     false,
		HideValue: true,
	})
	recreateWeb := cmd.Bool(cli.BoolOpt{
		Name:      "w web",
		Desc:      "Recreate existing web container",
		Value:     false,
		HideValue: true,
	})
	recreateDatabase := cmd.Bool(cli.BoolOpt{
		Name:      "d database",
		Desc:      "Recreate existing database container",
		Value:     false,
		HideValue: true,
	})

	cmd.Action = doRestartService(recreate, recreateServer, recreateWeb, recreateDatabase)
}

func doRestartService(recreate, recreateServer, recreateWeb, recreateDatabase *bool) func() {
	return func() {
		doStopService()()

		client, err := docker.NewDocker("")
		if err != nil {
			log.Fatal(err)
		}

		config, err := getConfig()
		if err != nil {
			log.Fatal(err)
		}

		allContainers, err := client.Ps(&docker.PsOptions{
			All: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		if *recreate {
			if !*recreateDatabase && !*recreateServer && !*recreateWeb {
				*recreateServer = true
				*recreateWeb = true
			}

			if *recreateServer {
				err := destroyContainer(client, bzkContainerServer, allContainers)
				if err != nil {
					log.Fatal(err)
				}
			}

			if *recreateWeb {
				err := destroyContainer(client, bzkContainerWeb, allContainers)
				if err != nil {
					log.Fatal(err)
				}
			}

			if *recreateDatabase {
				err := destroyContainer(client, bzkContainerMongo, allContainers)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		doStartService(&config.Tag, &config.Home, &config.DockerSock, &config.Registry, &config.SCMKey, &config.SyslogURL, &config.DbURL)()
	}
}

func statusService(cmd *cli.Cmd) {
	cmd.Action = func() {
		client, err := docker.NewDocker("")
		if err != nil {
			log.Fatal(err)
		}

		config, err := getConfig()
		if err != nil {
			log.Fatal(err)
		}

		allContainers, err := client.Ps(&docker.PsOptions{
			All: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		mongoUp := true
		if config.DbURL == "" {
			mongoUp = getContainerStatus(bzkContainerMongo, allContainers)
		}
		serverUp := getContainerStatus(bzkContainerServer, allContainers)
		webUp := getContainerStatus(bzkContainerWeb, allContainers)

		if mongoUp && serverUp && webUp {
			fmt.Printf("Bazooka service is Up\n")
		} else {
			fmt.Printf("Bazooka service is Down\n")
		}

	}
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

		// Recreate containers except database
		recreate, recreateServer, recreateWeb, recreateDatabase := true, true, true, false
		doRestartService(&recreate, &recreateServer, &recreateWeb, &recreateDatabase)()
	}
}

func getContainerStatus(name string, allContainers []dockerclient.APIContainers) bool {
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

func getConfigWithParams(tag, bzkHome, dockerSock, registry, scmKey, syslogURL, dbURL string) (*Config, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to load Bazooka config, reason is: %v\n", err)
	}
	currentUser, errCurrentUser := user.Current()

	if len(bzkHome) == 0 {
		if len(config.Home) == 0 {
			defaultHome := ""
			if errCurrentUser == nil {
				defaultHome = currentUser.HomeDir + "/bazooka"
			}

			config.Home = interactiveInput("Bazooka Home Folder", defaultHome)
		}
	} else {
		config.Home = bzkHome
	}

	if len(dockerSock) == 0 {
		if len(config.DockerSock) == 0 {
			config.DockerSock = interactiveInput("Docker Socket path", "/var/run/docker.sock")
		}
	} else {
		config.DockerSock = dockerSock
	}

	if len(scmKey) != 0 {
		config.SCMKey = scmKey
	}

	if len(*bzkApiUrl) != 0 {
		config.ApiURL = *bzkApiUrl
	}

	if len(syslogURL) != 0 {
		config.SyslogURL = syslogURL
	}

	if len(dbURL) != 0 {
		config.DbURL = dbURL
	}

	if len(registry) != 0 {
		config.Registry = registry
	}

	if len(tag) == 0 {
		if len(config.Tag) == 0 {
			config.Tag = "latest"
		}
	} else {
		config.Tag = tag
	}

	err = saveConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Unable to save Bazooka config, reason is: %v\n", err)
	}

	return config, nil
}

func getConfig() (*Config, error) {
	return getConfigWithParams("", "", "", "", "", "", "")
}

func destroyContainer(client *docker.Docker, name string, allContainers []dockerclient.APIContainers) error {
	container, err := getContainer(allContainers, name)
	if err != nil {
		fmt.Printf("Container %s does not exist\n", name)
		return nil
	}
	fmt.Printf("Destroying Container %s\n", name)
	err = client.Rm(&docker.RmOptions{
		Container: []string{container.ID},
		Force:     true,
	})
	return err
}

func stopContainer(client *docker.Docker, name string, allContainers []dockerclient.APIContainers) error {
	container, err := getContainer(allContainers, name)
	if err != nil {
		fmt.Printf("Container %s not found, doing nothing\n", name)
		return nil
	}
	if strings.HasPrefix(container.Status, "Up") {
		fmt.Printf("Stopping Container %s\n", name)
		err = client.Stop(&docker.StopOptions{
			ID:      container.ID,
			Timeout: 5,
		})
		if err != nil {
			return fmt.Errorf("Error stopping container %s, reason is %v\n", name, err)
		}
	} else {
		fmt.Printf("Container %s not running, doing nothing\n", name)
	}

	return nil
}

func startContainer(client *docker.Docker, options *docker.RunOptions, allContainers []dockerclient.APIContainers) error {
	container, err := getContainer(allContainers, options.Name)
	if err != nil {
		fmt.Printf("Container %s not found, creating and starting it\n", options.Name)
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
			ID:              container.ID,
			VolumeBinds:     options.VolumeBinds,
			Links:           options.Links,
			PublishAllPorts: options.PublishAllPorts,
			PortBindings:    options.PortBindings,
		})
	}

	fmt.Printf("Container %s found, but not in the right version, recreating it\n", options.Name)
	err = destroyContainer(client, options.Name, allContainers)
	if err != nil {
		return err
	}

	fmt.Printf("Starting %s\n", options.Name)
	_, err = client.Run(options)
	return err

}

func getServerEnv(home, dockerSock, scmKey, syslogURL, dbURL string) map[string]string {
	envMap := map[string]string{
		"BZK_HOME":       home,
		"BZK_DOCKERSOCK": dockerSock,
		"BZK_API_URL":    fmt.Sprintf("http://%s:3000", bzkContainerServer),
		"BZK_SYSLOG_URL": syslogURL,
		"BZK_NETWORK":    bzkNetwork,
	}
	if len(scmKey) > 0 {
		envMap["BZK_SCM_KEYFILE"] = scmKey
	}
	if len(dbURL) == 0 {
		dbURL = fmt.Sprintf("%s:27017", bzkContainerMongo)
	}
	envMap["BZK_DB_URL"] = dbURL
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
		Name: bzkContainerMongo,
		// Using the official mongo image from dockerhub, this may need a change later
		Image:       "mongo:3.0.2",
		Detach:      true,
		NetworkMode: bzkNetwork,
	}
}

func getServerRunOptions(config *Config) *docker.RunOptions {
	links := []string{}

	if config.DbURL == "" {
		links = append(links, fmt.Sprintf("%s:mongo", bzkContainerMongo))
	}

	return &docker.RunOptions{
		Name:   bzkContainerServer,
		Image:  getImageLocation(config.Registry, "bazooka/server", config.Tag),
		Detach: true,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", config.Home),
			fmt.Sprintf("%s:/var/run/docker.sock", config.DockerSock),
		},
		NetworkMode: bzkNetwork,
		Env:         getServerEnv(config.Home, config.DockerSock, config.SCMKey, config.SyslogURL, config.DbURL),
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"3000/tcp": {{HostPort: "3000"}},
			"3001/tcp": {{HostPort: "3001"}},
		},
	}
}

func getWebRunOptions(config *Config) *docker.RunOptions {
	return &docker.RunOptions{
		Name:        bzkContainerWeb,
		Image:       getImageLocation(config.Registry, "bazooka/web", config.Tag),
		Detach:      true,
		NetworkMode: bzkNetwork,
		Env:         map[string]string{"BZK_SERVER_HOST": bzkContainerServer},
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"80/tcp": {{HostPort: "8000"}},
		},
	}
}

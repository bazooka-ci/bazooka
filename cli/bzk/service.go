package main

import (
	"fmt"
	"log"
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
	bzkContainerQueue  = "bzk_queue"
)

func startService(cmd *cli.Cmd) {
	cmd.Spec = "[OPTIONS]"

	dbURL := cmd.String(cli.StringOpt{
		Name:   "db-url",
		Desc:   "URL of a MongoDB server",
		EnvVar: "BZK_DB_URL",
	})
	queueURL := cmd.String(cli.StringOpt{
		Name:   "queue-url",
		Desc:   "URL of the job queue",
		EnvVar: "BZK_QUEUE_URL",
	})
	registry := cmd.String(cli.StringOpt{
		Name:   "registry",
		EnvVar: "BZK_REGISTRY",
	})
	tag := cmd.String(cli.StringOpt{
		Name: "tag",
		Desc: "The bazooka version to run",
	})

	cmd.Action = doStartService(tag, registry, dbURL, queueURL)
}

func doStartService(version, registry, dbURL, queueURL *string) func() {
	return func() {
		ensureNetworkExists()

		config, err := getConfigWithParams(*version, *registry, *dbURL, *queueURL)
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

		if config.QueueURL == "" {
			err = startContainer(client, getQueueRunOptions(), allContainers)
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
	cmd.Spec = "[-x [-swdq]]"

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
	recreateQueue := cmd.Bool(cli.BoolOpt{
		Name:      "q queue",
		Desc:      "Recreate existing queue container",
		Value:     false,
		HideValue: true,
	})

	cmd.Action = doRestartService(recreate, recreateServer, recreateWeb, recreateDatabase, recreateQueue)
}

func doRestartService(recreate, recreateServer, recreateWeb, recreateDatabase, recreateQueue *bool) func() {
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
			if !*recreateQueue && !*recreateDatabase && !*recreateServer && !*recreateWeb {
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

			if *recreateQueue {
				err := destroyContainer(client, bzkContainerQueue, allContainers)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		doStartService(&config.Tag, &config.Registry, &config.DbURL, &config.QueueURL)()
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
		queueUp := true
		if config.QueueURL == "" {
			queueUp = getContainerStatus(bzkContainerQueue, allContainers)
		}
		serverUp := getContainerStatus(bzkContainerServer, allContainers)
		webUp := getContainerStatus(bzkContainerWeb, allContainers)

		if queueUp && mongoUp && serverUp && webUp {
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

		// Recreate containers except database and queue
		recreate, recreateServer, recreateWeb, recreateDatabase, recreateQueue := true, true, true, false, false
		doRestartService(&recreate, &recreateServer, &recreateWeb, &recreateDatabase, &recreateQueue)()
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

func getConfigWithParams(tag, registry, dbURL, queueURL string) (*Config, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to load Bazooka config, reason is: %v\n", err)
	}

	if len(dbURL) != 0 {
		config.DbURL = dbURL
	}

	if len(queueURL) != 0 {
		config.QueueURL = queueURL
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
	return getConfigWithParams("", "", "", "")
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

func getServerEnv(dbURL, queueURL string) map[string]string {
	envMap := map[string]string{
		"BZK_QUEUE_URL": queueURL,
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

func getQueueRunOptions() *docker.RunOptions {
	return &docker.RunOptions{
		Name:   bzkContainerQueue,
		Image:  "jawher/beanstalkd",
		Detach: true,
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"11300/tcp": {{HostPort: "11300"}},
		},
	}
}

func getServerRunOptions(config *Config) *docker.RunOptions {
	links := []string{}

	if config.DbURL == "" {
		links = append(links, fmt.Sprintf("%s:mongo", bzkContainerMongo))
	}

	return &docker.RunOptions{
		Name:        bzkContainerServer,
		Image:       getImageLocation(config.Registry, "bazooka/server", config.Tag),
		Detach:      true,
		NetworkMode: bzkNetwork,
		Env:         getServerEnv(config.DbURL, config.QueueURL),
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

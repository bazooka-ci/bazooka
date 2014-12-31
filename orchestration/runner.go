package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	docker "github.com/bywan/go-dockercommand"
	commons "github.com/haklop/bazooka/commons"
	"github.com/haklop/bazooka/commons/mongo"
	"github.com/haklop/bazooka/commons/parallel"
)

type Runner struct {
	Variants []*variantData
	Env      map[string]string
	Mongo    *mongo.MongoConnector
	client   *docker.Docker
}

func (r *Runner) Run(logger Logger) error {
	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return err
	}
	r.client = client

	par := parallel.New()

	for _, ivariant := range r.Variants {
		if ivariant.variant.Status != commons.JOB_RUNNING {
			continue
		}
		variant := ivariant
		par.Submit(func() error {
			return r.runContainer(logger, variant, r.Env)
		}, variant)
	}

	par.Exec(func(tag interface{}, err error) {
		v := tag.(*variantData)
		if err != nil {
			log.Errorf("Run error %v for variant %v\n", err, v)
			v.variant.Status = commons.JOB_ERRORED
		} else {
			log.WithFields(log.Fields{
				"variant": v.counter,
			}).Info("Variant Completed")
		}
		v.variant.Completed = time.Now()
	})

	log.Info("Dockerfiles builds finished")
	return nil
}

func (r *Runner) runContainer(logger Logger, vd *variantData, env map[string]string) error {
	success := true
	servicesFile := fmt.Sprintf("%s/work/%d/services", BazookaInput, vd.counter)

	servicesList, err := listServices(servicesFile)
	if err != nil {
		return err
	}

	serviceContainers := []*docker.Container{}
	containerLinks := []string{}
	for _, service := range servicesList {
		name := fmt.Sprintf("service-%s-%s-%d", env[BazookaEnvProjectID], env[BazookaEnvJobID], vd.variant.Number)
		containerLinks = append(containerLinks, fmt.Sprintf("%s:%s", name, service))
		serviceContainer, err := r.client.Run(&docker.RunOptions{
			Name:   name,
			Image:  service,
			Detach: true,
		})
		if err != nil {
			return err
		}
		serviceContainers = append(serviceContainers, serviceContainer)
	}

	// TODO link containers
	container, err := r.client.Run(&docker.RunOptions{
		Image: vd.imageTag,
		Links: containerLinks,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/var/run/docker.sock", DockerSock),
		},
		Detach: true,
	})
	if err != nil {
		return err
	}

	container.Logs(vd.imageTag)
	logger(vd.imageTag, vd.variant.ID, container)

	exitCode, err := container.Wait()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		if exitCode == 42 {
			return fmt.Errorf("Run failed\n Check Docker container logs, id is %s\n", container.ID())
		}
		success = false
	}
	if err = container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {
		return err
	}

	for _, serviceContainer := range serviceContainers {
		if err = serviceContainer.Remove(&docker.RemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		}); err != nil {
			return err
		}
	}

	if success {
		vd.variant.Status = commons.JOB_SUCCESS
	} else {
		vd.variant.Status = commons.JOB_FAILED
	}
	return nil
}

func listServices(servicesFile string) ([]string, error) {
	file, err := os.Open(servicesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var services []string
	for scanner.Scan() {
		services = append(services, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

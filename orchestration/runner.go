package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"

	commons "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/commons/mongo"
	"github.com/bazooka-ci/bazooka/commons/parallel"
	docker "github.com/bywan/go-dockercommand"
)

var (
	whiteListEnvVarsNames = []string{
		"BZK_SCM_URL", "BZK_PROJECT_ID", "BZK_JOB_ID", "BZK_BUILD_DIR",
		"BZK_JOB_PARAMETERS", "BZK_SCM", "BZK_SCM_REFERENCE", "BZK_VARIANT",
	}
)

type Runner struct {
	Variants            []*variantData
	ArtifactsFolderBase string
	Env                 map[string]string
	Mongo               *mongo.MongoConnector
	client              *docker.Docker
}

func (r *Runner) Run(logger Logger) error {
	client, err := docker.NewDocker(paths.container.dockerEndpoint)
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
			return r.runContainer(logger, variant)
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

func (r *Runner) runContainer(logger Logger, vd *variantData) error {
	success := true
	servicesFile := fmt.Sprintf("%s/%d/services", paths.container.work, vd.counter)

	servicesList, err := listServices(servicesFile)
	if err != nil {
		return err
	}

	serviceContainers := []*docker.Container{}
	containerLinks := []string{}
	for _, service := range servicesList {
		name := fmt.Sprintf("service-%s-%s-%d", r.Env[BazookaEnvProjectID], r.Env[BazookaEnvJobID], vd.variant.Number)
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

	artifactsFolder := fmt.Sprintf("%s/%s", r.ArtifactsFolderBase, vd.variant.ID)

	container, err := r.client.Run(&docker.RunOptions{
		Image: vd.imageTag,
		Links: containerLinks,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/var/run/docker.sock", paths.host.dockerSock),
			fmt.Sprintf("%s:/artifacts", artifactsFolder),
		},
		Env:    whiteListEnvVars(injectVariantInEnv(vd.variant.Number, r.Env)),
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

func whiteListEnvVars(envVars map[string]string) map[string]string {
	res := map[string]string{}
	for k, v := range envVars {
		if stringInSlice(k, whiteListEnvVarsNames) {
			res[k] = v
		}
	}
	return res
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func injectVariantInEnv(v int, env map[string]string) map[string]string {
	env["BZK_VARIANT"] = strconv.Itoa(v)
	return env
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

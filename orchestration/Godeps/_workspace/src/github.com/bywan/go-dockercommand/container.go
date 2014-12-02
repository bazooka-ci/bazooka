package dockercommand

import (
	"bufio"
	"io"
	"log"

	docker "github.com/fsouza/go-dockerclient"
)

type Container struct {
	info   *docker.Container
	client *docker.Client
}

type RemoveOptions struct {
	RemoveVolumes bool
	Force         bool
}

func (c *Container) Wait() (int, error) {
	return c.client.WaitContainer(c.info.ID)
}

func (c *Container) Stop(timeout uint) error {
	return c.client.StopContainer(c.info.ID, timeout)
}

func (c *Container) StreamLogs(w io.Writer) {
	options := docker.LogsOptions{
		Container:    c.info.ID,
		OutputStream: w,
		ErrorStream:  w,
		Follow:       true,
		Stdout:       true,
		Stderr:       true,
	}
	go func() {
		err := c.client.Logs(options)
		if err != nil {
			log.Println(err.Error())
		}
	}()
}

func (c *Container) Logs(prefix string) {
	r, w := io.Pipe()
	c.StreamLogs(w)
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Printf("[%s] %s \n", prefix, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Println("There was an error with the scanner in attached container", err)
		}
	}(r)
}

func (c *Container) Remove(opts *RemoveOptions) error {
	options := docker.RemoveContainerOptions{
		ID:            c.info.ID,
		Force:         opts.Force,
		RemoveVolumes: opts.RemoveVolumes,
	}
	return c.client.RemoveContainer(options)
}

func (c *Container) ID() string {
	return c.info.ID
}

func (c *Container) Inspect() (*docker.Container, error) {
	return c.client.InspectContainer(c.info.ID)
}

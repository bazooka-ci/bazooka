package dockercommand

import docker "github.com/fsouza/go-dockerclient"

type RmOptions struct {
	Container     []string
	Force         bool
	RemoveVolumes bool
}

func (dock *Docker) Rm(options *RmOptions) error {
	for _, containerID := range options.Container {
		err := dock.client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            containerID,
			Force:         options.Force,
			RemoveVolumes: options.RemoveVolumes,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

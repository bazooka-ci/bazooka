package dockercommand

import docker "github.com/fsouza/go-dockerclient"

func (dock *Docker) Inspect(containerID string) (*docker.Container, error) {
	container, err := dock.client.InspectContainer(containerID)
	if err != nil {
		return nil, err
	}

	return container, nil
}

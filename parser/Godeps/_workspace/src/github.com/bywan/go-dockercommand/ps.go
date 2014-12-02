package dockercommand

import docker "github.com/fsouza/go-dockerclient"

type PsOptions struct {
	All    bool
	Size   bool
	Limit  int
	Since  string
	Before string
}

func (dock *Docker) Ps(options *PsOptions) ([]docker.APIContainers, error) {
	containers, err := dock.client.ListContainers(docker.ListContainersOptions{
		All:    options.All,
		Size:   options.Size,
		Limit:  options.Limit,
		Since:  options.Since,
		Before: options.Before,
	})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

package dockercommand

import (
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

func (dock *Docker) pullImageIfNotExist(image string) error {
	_, err := dock.client.InspectImage(image)
	fmt.Printf("1\n")
	if err != nil && err.Error() == "no such image" {
		fmt.Printf("2\n")
		err = dock.client.PullImage(docker.PullImageOptions{Repository: image, OutputStream: os.Stdout},
			docker.AuthConfiguration{})
	}
	return err
}

func convertEnvMapToSlice(envMap map[string]string) []string {
	envSlice := []string{}
	for key, value := range envMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", key, value))
	}
	return envSlice
}

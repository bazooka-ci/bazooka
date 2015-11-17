package bazooka

import (
	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
)

// RemoveContainer will cleanly remove a Docker container and warn user of a leak in case of error
func RemoveContainer(container *docker.Container) {
	err := container.Remove(&docker.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"name":  container.ID(),
			"error": err.Error(),
		}).Error("Error while removing service container, Be aware that this could cause container leaks if services have not been removed correctly")
	}
}

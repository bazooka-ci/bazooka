package main

import (
	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
)

func (c *context) addKey(params map[string]string, body bodyFunc) (*response, error) {
	var key lib.SSHKey

	body(&key)

	if len(key.Content) == 0 {
		return badRequest("content is mandatory")
	}

	project, err := c.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	keys, err := c.Connector.GetKeys(params["id"])
	if err != nil {
		return nil, err
	}

	if len(keys) > 0 {
		return conflict("A key is already associated with this project")
	}

	key.ProjectID = project.ID

	log.WithFields(log.Fields{
		"key": key,
	}).Debug("Adding key")

	if err = c.Connector.AddKey(&key); err != nil {
		return nil, err
	}

	createdKey := &lib.SSHKey{
		ProjectID: key.ProjectID,
	}

	return created(&createdKey, "/project/"+params["id"]+"/key")
}

func (c *context) updateKey(params map[string]string, body bodyFunc) (*response, error) {
	var key lib.SSHKey

	body(&key)

	if len(key.Content) == 0 {
		return badRequest("content is mandatory")
	}

	project, err := c.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	key.ProjectID = project.ID

	log.WithFields(log.Fields{
		"key": key,
	}).Debug("Updating key")

	if err = c.Connector.UpdateKey(project.ID, &key); err != nil {
		return nil, err
	}

	updateKey := &lib.SSHKey{
		ProjectID: key.ProjectID,
	}

	return ok(&updateKey)
}

func (c *context) listKeys(params map[string]string, body bodyFunc) (*response, error) {

	keys, err := c.Connector.GetKeys(params["id"])
	if err != nil {
		return nil, err
	}

	return ok(&keys)

}

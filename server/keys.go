package main

import (
	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
)

func (c *context) addKey(r *request) (*response, error) {
	var key lib.SSHKey

	r.parseBody(&key)

	if len(key.Content) == 0 {
		return badRequest("content is mandatory")
	}

	project, err := c.connector.GetProjectById(r.vars["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	keys, err := c.connector.GetKeys(r.vars["id"])
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

	if err = c.connector.AddKey(&key); err != nil {
		return nil, err
	}

	createdKey := &lib.SSHKey{
		ProjectID: key.ProjectID,
	}

	return created(&createdKey, "/project/"+r.vars["id"]+"/key")
}

func (c *context) updateKey(r *request) (*response, error) {
	var key lib.SSHKey

	r.parseBody(&key)

	if len(key.Content) == 0 {
		return badRequest("content is mandatory")
	}

	project, err := c.connector.GetProjectById(r.vars["id"])
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

	if err = c.connector.UpdateKey(project.ID, &key); err != nil {
		return nil, err
	}

	updateKey := &lib.SSHKey{
		ProjectID: key.ProjectID,
	}

	return ok(&updateKey)
}

func (c *context) listKeys(r *request) (*response, error) {

	keys, err := c.connector.GetKeys(r.vars["id"])
	if err != nil {
		return nil, err
	}

	return ok(&keys)
}

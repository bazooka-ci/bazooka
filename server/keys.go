package main

import (
	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/server/db"
)

func (c *context) setKey(r *request) (*response, error) {
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
	}).Debug("Setting key")

	if err = c.connector.SetProjectKey(project.ID, &key); err != nil {
		return nil, err
	}

	return noContent()
}

func (c *context) getKey(r *request) (*response, error) {
	key, err := c.connector.GetProjectKey(r.vars["id"])
	if err != nil {
		if _, ok := err.(*db.NotFoundError); ok {
			return notFound("key not found")
		}
		return nil, err
	}

	return ok(&key)
}

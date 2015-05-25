package db

import (
	lib "github.com/bazooka-ci/bazooka/commons"
	"gopkg.in/mgo.v2/bson"
)

func (c *MongoConnector) GetProjectKey(projectID string) (*lib.SSHKey, error) {
	proj, err := c.GetProjectById(projectID)
	if err != nil {
		return nil, err
	}

	result := &lib.SSHKey{}
	if err := c.selectOneByField("keys", "project_id", proj.ID, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) SetProjectKey(projectID string, key *lib.SSHKey) error {
	proj, err := c.GetProjectById(projectID)
	if err != nil {
		return err
	}

	_, err = c.database.C("keys").Upsert(bson.M{
		"project_id": proj.ID,
	}, key)
	return err
}

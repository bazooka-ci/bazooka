package mongo

import (
	lib "github.com/haklop/bazooka/commons"
	"gopkg.in/mgo.v2/bson"
)

func (c *MongoConnector) GetProjectKey(projectID string) (*lib.SSHKey, error) {
	result := &lib.SSHKey{}
	if err := c.ByField("keys", "project_id", projectID, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) UpdateKey(id string, key *lib.SSHKey) error {
	return c.database.C("keys").Update(bson.M{
		"project_id": id,
	}, key)
}

func (c *MongoConnector) AddKey(key *lib.SSHKey) error {
	var err error
	if key.ID, err = c.randomId(); err != nil {
		return err
	}

	return c.database.C("keys").Insert(key)
}

func (c *MongoConnector) GetKeys(projectID string) ([]*lib.SSHKey, error) {
	proj, err := c.GetProjectById(projectID)
	if err != nil {
		return nil, err
	}

	result := []*lib.SSHKey{}

	err = c.database.C("keys").Find(bson.M{
		"project_id": proj.ID,
	}).All(&result)

	return result, err
}

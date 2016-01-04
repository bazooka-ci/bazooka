package db

import (
	"github.com/bazooka-ci/bazooka/commons"
	"gopkg.in/mgo.v2/bson"
)

func (c *MongoConnector) GetProjectCryptoKey(projectID string) (*bazooka.CryptoKey, error) {
	result := &bazooka.CryptoKey{}
	if err := c.selectOneByField("crypto", "project_id", projectID, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) UpdateCryptoKey(id string, key *bazooka.CryptoKey) error {
	return c.database.C("crypto").Update(bson.M{
		"project_id": id,
	}, key)
}

func (c *MongoConnector) HasCryptoKey(projectID string) (bool, error) {
	request := bson.M{"project_id": projectID}
	count, err := c.database.C("crypto").Find(request).Count()
	return count > 0, err
}

func (c *MongoConnector) AddCryptoKey(key *bazooka.CryptoKey) error {
	hasKey, err := c.HasCryptoKey(key.ProjectID)
	if err != nil {
		return err
	}
	if hasKey {
		return c.UpdateCryptoKey(key.ProjectID, key)
	}
	if key.ID, err = c.randomId(); err != nil {
		return err
	}

	return c.database.C("crypto").Insert(key)
}

func (c *MongoConnector) GetCryptoKeys(projectID string) ([]*bazooka.CryptoKey, error) {
	proj, err := c.GetProjectById(projectID)
	if err != nil {
		return nil, err
	}

	result := []*bazooka.CryptoKey{}

	err = c.database.C("crypto").Find(bson.M{
		"project_id": proj.ID,
	}).All(&result)

	return result, err
}

package db

import (
	lib "github.com/bazooka-ci/bazooka/commons"

	"gopkg.in/mgo.v2/bson"
)

func (c *MongoConnector) HasImage(name string) (bool, error) {
	request := bson.M{"name": name}
	count, err := c.database.C("images").Find(request).Count()
	return count > 0, err
}

func (c *MongoConnector) GetImage(name string) (*lib.Image, error) {
	request := bson.M{"name": name}
	im := lib.Image{}
	err := c.database.C("images").Find(request).One(&im)
	return &im, err
}

func (c *MongoConnector) SetImage(name, image string) error {
	selector := bson.M{
		"name": name,
	}
	request := bson.M{
		"$set": bson.M{"image": image},
	}
	_, err := c.database.C("images").Upsert(selector, request)
	return err
}

func (c *MongoConnector) GetImages() ([]*lib.Image, error) {

	res := []*lib.Image{}
	err := c.database.C("images").Find(bson.M{}).All(&res)
	return res, err
}

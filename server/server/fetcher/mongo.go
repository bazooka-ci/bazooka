package fetcher

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName = "fetchers"
)

type mongoConnector struct {
	Database *mgo.Database
}

func (c *mongoConnector) GetFetcherByName(name string) (Fetcher, error) {
	result := Fetcher{}

	request := bson.M{
		"name": name,
	}
	err := c.Database.C(collectionName).Find(request).One(&result)
	return result, err
}

func (c *mongoConnector) GetFetcherById(id string) (Fetcher, error) {
	result := Fetcher{}

	request := bson.M{
		"id": id,
	}
	err := c.Database.C(collectionName).Find(request).One(&result)
	return result, err
}

func (c *mongoConnector) GetFetchers() ([]Fetcher, error) {
	result := []Fetcher{}

	err := c.Database.C(collectionName).Find(bson.M{}).All(&result)
	return result, err
}

func (c *mongoConnector) AddFetcher(fetcher *Fetcher) error {
	i := bson.NewObjectId()
	fetcher.ID = i.Hex()

	err := c.Database.C(collectionName).Insert(fetcher)
	return err
}

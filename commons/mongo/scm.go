package mongo

import (
	lib "github.com/haklop/bazooka/commons"
	"gopkg.in/mgo.v2/bson"
)

const (
	scmCollectionName = "fetchers"
)

func (c *MongoConnector) GetFetcherByName(name string) (lib.ScmFetcher, error) {
	result := lib.ScmFetcher{}

	request := bson.M{
		"name": name,
	}
	err := c.database.C(scmCollectionName).Find(request).One(&result)
	return result, err
}

func (c *MongoConnector) GetFetcherById(id string) (lib.ScmFetcher, error) {
	result := lib.ScmFetcher{}

	request := bson.M{
		"id": id,
	}
	err := c.database.C(scmCollectionName).Find(request).One(&result)
	return result, err
}

func (c *MongoConnector) GetFetchers() ([]lib.ScmFetcher, error) {
	result := []lib.ScmFetcher{}

	err := c.database.C(scmCollectionName).Find(bson.M{}).All(&result)
	return result, err
}

func (c *MongoConnector) AddFetcher(fetcher *lib.ScmFetcher) error {
	i := bson.NewObjectId()
	fetcher.ID = i.Hex()

	err := c.database.C(scmCollectionName).Insert(fetcher)
	return err
}

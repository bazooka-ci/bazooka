package mongo

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

const (
	bazookaEnvMongoAddr = "MONGO_PORT_27017_TCP_ADDR"
	bazookaEnvMongoPort = "MONGO_PORT_27017_TCP_PORT"
	bazookaMongoBase    = "bazooka"
)

type MongoConnector struct {
	database *mgo.Database
	session  *mgo.Session
}

func NewConnector() *MongoConnector {
	session, err := mgo.Dial(os.Getenv(bazookaEnvMongoAddr) + ":" + os.Getenv(bazookaEnvMongoPort))
	if err != nil {
		panic(err)
	}

	database := session.DB(bazookaMongoBase)

	connector := &MongoConnector{
		database: database,
		session:  session,
	}

	return connector
}

func (c *MongoConnector) Close() {
	c.session.Close()
}

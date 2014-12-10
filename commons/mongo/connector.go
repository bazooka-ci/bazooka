package mongo

import (
	"crypto/rand"
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

type NotFoundError struct {
	Collection string
	Field      string
	Value      string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("%s[%s:%s] not found", n.Collection, n.Field, n.Value)
}

type ManyFoundError struct {
	Collection string
	Field      string
	Value      string
	Count      int
}

func (m *ManyFoundError) Error() string {
	return fmt.Sprintf("%s[%s:%s] returned %d results", m.Collection, m.Field, m.Value, m.Count)
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

func (m *MongoConnector) idLike(id string) bson.M {
	return bson.M{
		"id": bson.M{
			"$regex":   "^" + id,
			"$options": "i",
		},
	}
}

func (m *MongoConnector) randomId() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func (m *MongoConnector) ByField(collection, fieldName string, fieldValue string, result interface{}) error {
	q := m.database.C(collection).Find(bson.M{fieldName: fieldValue})
	count, err := q.Count()
	if err != nil {
		return err
	}
	switch count {
	case 0:
		return &NotFoundError{collection, fieldName, fieldValue}
	case 1:
		return q.One(result)

	default:
		return &ManyFoundError{collection, fieldName, fieldValue, count}
	}
}

func (m *MongoConnector) ById(collection, id string, result interface{}) error {
	q := m.database.C(collection).Find(m.idLike(id))
	count, err := q.Count()
	if err != nil {
		return err
	}
	switch count {
	case 0:
		return &NotFoundError{collection, "id", id}
	case 1:
		return q.One(result)

	default:
		return &ManyFoundError{collection, "id", id, count}
	}
}

func (m *MongoConnector) ByIdOrName(collection, id string, result interface{}) error {
	err := m.ById(collection, id, result)
	switch err.(type) {
	case *NotFoundError:
		return m.ByField(collection, "name", id, result)
	default:
		return err
	}

}

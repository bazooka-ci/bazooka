package db

import (
	"crypto/rand"
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	bazookaMongoBase = "bazooka"
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

func NewConnector(url string) *MongoConnector {
	session, err := mgo.Dial(url)
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

func (m *MongoConnector) fieldStartsWith(field, value string) bson.M {
	return bson.M{
		field: bson.M{
			"$regex":   "^" + value,
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

func (m *MongoConnector) selectOneByField(collection, fieldName string, fieldValue string, result interface{}) error {
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

func (m *MongoConnector) selectOneByFieldLike(collection, fieldName string, fieldValue string, result interface{}) error {
	q := m.database.C(collection).Find(m.fieldStartsWith(fieldName, fieldValue))
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

func (m *MongoConnector) selectOneByIdOrName(collection, id string, result interface{}) error {
	err := m.selectOneByFieldLike(collection, "id", id, result)
	switch err.(type) {
	case *NotFoundError:
		return m.selectOneByField(collection, "name", id, result)
	default:
		return err
	}
}

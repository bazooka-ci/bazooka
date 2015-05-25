package db

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	lib "github.com/bazooka-ci/bazooka/commons"
)

func (c *MongoConnector) GetUserByEmail(email string) (*lib.User, error) {
	result := &lib.User{}
	if err := c.selectOneByField("user", "email", email, result); err != nil {
		return nil, err
	}
	result.Password = ""
	fmt.Printf("retrieve user: %#v\n", result)
	return result, nil
}

func (c *MongoConnector) HasUser(email string) (bool, error) {
	request := bson.M{}
	request["email"] = email

	count, err := c.database.C("user").Find(request).Count()
	return count > 0, err
}

func (c *MongoConnector) GetUsers() ([]*lib.User, error) {
	result := []*lib.User{}

	err := c.database.C("user").Find(bson.M{}).All(&result)
	for _, user := range result {
		user.Password = ""
	}
	fmt.Printf("retrieve users: %#v\n", result)
	return result, err
}

func (c *MongoConnector) AddUser(user *lib.User) error {
	var err error
	if user.ID, err = c.randomId(); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 2) // TODO define a smart value
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	fmt.Printf("add user: %s\n", user.Email)
	err = c.database.C("user").Insert(user)
	if err != nil {
		return err
	} else {
		user.Password = ""
		return nil
	}
}

func (c *MongoConnector) ComparePassword(email string, password string) bool {
	result := &lib.User{}
	if err := c.selectOneByField("user", "email", email, result); err != nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	return err == nil
}

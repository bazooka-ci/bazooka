package client

import (
	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

type User struct {
	config *Config
}

func (c *User) List() ([]lib.User, error) {
	var u []lib.User

	requestURL, err := c.config.getRequestURL("user")
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &u,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return u, err
}

func (c *User) Create(email, password string) (*lib.User, error) {
	user := lib.User{
		Email:    email,
		Password: password,
	}
	createdUser := &lib.User{}

	requestURL, err := c.config.getRequestURL("user")
	if err != nil {
		return nil, err
	}

	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &user,
		Results:    &createdUser,
		OkCodes:    []int{201},
		SetHeaders: c.config.authenticateRequest,
	})

	return createdUser, err
}

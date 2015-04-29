package client

import (
	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

func (c *Client) ListUsers() ([]lib.User, error) {
	var u []lib.User

	requestURL, err := c.getRequestURL("user")
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &u,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return u, err
}

func (c *Client) CreateUser(email, password string) (*lib.User, error) {
	user := lib.User{
		Email:    email,
		Password: password,
	}
	createdUser := &lib.User{}

	requestURL, err := c.getRequestURL("user")
	if err != nil {
		return nil, err
	}

	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &user,
		Results:    &createdUser,
		OkCodes:    []int{201},
		SetHeaders: c.authenticateRequest,
	})

	return createdUser, err
}

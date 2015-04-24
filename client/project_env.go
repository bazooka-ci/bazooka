package client

import (
	"fmt"
	"strings"

	"github.com/racker/perigee"
)

type ProjectEnv struct {
	config *Config
}

func (c *ProjectEnv) Get(id string) (map[string]string, error) {
	var res map[string]string

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/env", id))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &res,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})
	return res, err
}

func (c *ProjectEnv) SetEnv(id, key, value string) error {
	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/env/%s", id, key))
	if err != nil {
		return err
	}
	return perigee.Put(requestURL, perigee.Options{
		ContentType: "text/plain",
		ReqBody:     strings.NewReader(value),
		OkCodes:     []int{204},
		SetHeaders:  c.config.authenticateRequest,
	})
}

func (c *ProjectEnv) UnsetEnv(id, key string) error {
	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/env/%s", id, key))
	if err != nil {
		return err
	}
	return perigee.Delete(requestURL, perigee.Options{
		OkCodes:    []int{204},
		SetHeaders: c.config.authenticateRequest,
	})
}

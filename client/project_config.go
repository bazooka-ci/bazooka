package client

import (
	"fmt"
	"strings"

	"github.com/racker/perigee"
)

type ProjectConfig struct {
	config *Config
}

func (c *ProjectConfig) Get(id string) (map[string]string, error) {
	var res map[string]string

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/config", id))
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

func (c *ProjectConfig) SetKey(id, key, value string) error {
	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/config/%s", id, key))
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

func (c *ProjectConfig) UnsetKey(id, key string) error {
	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/config/%s", id, key))
	if err != nil {
		return err
	}
	return perigee.Delete(requestURL, perigee.Options{
		OkCodes:    []int{204},
		SetHeaders: c.config.authenticateRequest,
	})
}

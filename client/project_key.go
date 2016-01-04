package client

import (
	"fmt"
	"io/ioutil"
	"net/url"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

type ProjectKey struct {
	config *Config
}

func (c *ProjectKey) Get(projectID string) (*lib.SSHKey, error) {
	var key lib.SSHKey

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &key,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return &key, err
}

func (c *ProjectKey) Set(projectID, keyPath string) error {
	fileContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}

	sshKey := &lib.SSHKey{
		ProjectID: projectID,
		Content:   string(fileContent),
	}

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return err
	}

	err = perigee.Put(requestURL, perigee.Options{
		ReqBody:    &sshKey,
		OkCodes:    []int{204},
		SetHeaders: c.config.authenticateRequest,
	})

	return err
}

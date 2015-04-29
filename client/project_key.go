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

func (c *ProjectKey) List(projectID string) ([]*lib.SSHKey, error) {

	var keys []*lib.SSHKey

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &keys,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return keys, err
}

func (c *ProjectKey) Add(projectID, keyPath string) (*lib.SSHKey, error) {
	fileContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	sshKey := &lib.SSHKey{
		ProjectID: projectID,
		Content:   string(fileContent),
	}

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	createdKey := &lib.SSHKey{}

	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &sshKey,
		Results:    &createdKey,
		OkCodes:    []int{201},
		SetHeaders: c.config.authenticateRequest,
	})

	return createdKey, err
}

func (c *ProjectKey) Update(projectID, keyPath string) (*lib.SSHKey, error) {
	fileContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	sshKey := &lib.SSHKey{
		ProjectID: projectID,
		Content:   string(fileContent),
	}

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	updatedKey := &lib.SSHKey{}

	err = perigee.Put(requestURL, perigee.Options{
		ReqBody:    &sshKey,
		Results:    &updatedKey,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return updatedKey, err
}

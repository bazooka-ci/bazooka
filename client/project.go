package client

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

func (c *Client) ListProjects() ([]lib.Project, error) {
	var p []lib.Project

	requestURL, err := c.getRequestURL("project")
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &p,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return p, err
}

func (c *Client) GetProjectConfig(id string) (map[string]string, error) {
	var res map[string]string

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/config", id))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &res,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})
	return res, err
}

func (c *Client) SetProjectConfigKey(id, key, value string) error {
	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/config/%s", id, key))
	if err != nil {
		return err
	}
	return perigee.Put(requestURL, perigee.Options{
		ContentType: "text/plain",
		ReqBody:     strings.NewReader(value),
		OkCodes:     []int{204},
		SetHeaders:  c.authenticateRequest,
	})
}

func (c *Client) UnsetProjectConfigKey(id, key string) error {

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/config/%s", id, key))
	if err != nil {
		return err
	}
	return perigee.Delete(requestURL, perigee.Options{
		OkCodes:    []int{204},
		SetHeaders: c.authenticateRequest,
	})
}

func (c *Client) ListJobs(projectID string) ([]lib.Job, error) {
	var j []lib.Job

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/job", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &j,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return j, err
}

func (c *Client) CreateProject(name, scm, scmUri string) (*lib.Project, error) {
	project := lib.Project{
		Name:    name,
		ScmType: scm,
		ScmURI:  scmUri,
	}
	createdProject := &lib.Project{}

	requestURL, err := c.getRequestURL("project")
	if err != nil {
		return nil, err
	}
	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &project,
		Results:    &createdProject,
		OkCodes:    []int{201},
		SetHeaders: c.authenticateRequest,
	})

	return createdProject, err
}

func (c *Client) StartJob(projectID, scmReference string, envParameters []string) (*lib.Job, error) {
	startJob := lib.StartJob{
		ScmReference: scmReference,
		Parameters:   envParameters,
	}
	createdJob := &lib.Job{}

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/job", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &startJob,
		Results:    &createdJob,
		OkCodes:    []int{202},
		SetHeaders: c.authenticateRequest,
	})

	return createdJob, err
}

func (c *Client) ListKeys(projectID string) ([]*lib.SSHKey, error) {

	var keys []*lib.SSHKey

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &keys,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return keys, err
}

func (c *Client) AddKey(projectID, keyPath string) (*lib.SSHKey, error) {
	fileContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	sshKey := &lib.SSHKey{
		ProjectID: projectID,
		Content:   string(fileContent),
	}

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	createdKey := &lib.SSHKey{}

	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &sshKey,
		Results:    &createdKey,
		OkCodes:    []int{201},
		SetHeaders: c.authenticateRequest,
	})

	return createdKey, err
}

func (c *Client) UpdateKey(projectID, keyPath string) (*lib.SSHKey, error) {
	fileContent, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	sshKey := &lib.SSHKey{
		ProjectID: projectID,
		Content:   string(fileContent),
	}

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/key", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	updatedKey := &lib.SSHKey{}

	err = perigee.Put(requestURL, perigee.Options{
		ReqBody:    &sshKey,
		Results:    &updatedKey,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return updatedKey, err
}

func (c *Client) EncryptData(projectID, toEncryptString string) (string, error) {
	toEncryptData := &lib.StringValue{
		Value: toEncryptString,
	}

	requestURL, err := c.getRequestURL(fmt.Sprintf("project/%s/crypto", url.QueryEscape(projectID)))
	if err != nil {
		return "", err
	}

	encryptedData := &lib.StringValue{}

	err = perigee.Put(requestURL, perigee.Options{
		ReqBody:    &toEncryptData,
		Results:    &encryptedData,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return encryptedData.Value, err
}

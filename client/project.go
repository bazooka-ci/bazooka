package client

import (
	"fmt"
	"net/url"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

type Project struct {
	config *Config
	Key    *ProjectKey
	Config *ProjectConfig
	Env    *ProjectEnv
}

func (c *Project) List() ([]lib.Project, error) {
	var p []lib.Project

	requestURL, err := c.config.getRequestURL("project")
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &p,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return p, err
}

func (c *Project) Jobs(projectID string) ([]lib.Job, error) {
	var j []lib.Job

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/job", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &j,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return j, err
}

func (c *Project) Get(projectID string) (*lib.Project, error) {
	var p lib.Project

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &p,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return &p, err
}

func (c *Project) Create(name, scm, scmUri string) (*lib.Project, error) {
	project := lib.Project{
		Name:    name,
		ScmType: scm,
		ScmURI:  scmUri,
	}
	createdProject := &lib.Project{}

	requestURL, err := c.config.getRequestURL("project")
	if err != nil {
		return nil, err
	}
	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &project,
		Results:    &createdProject,
		OkCodes:    []int{201},
		SetHeaders: c.config.authenticateRequest,
	})

	return createdProject, err
}

func (c *Project) StartJob(projectID, scmReference string, envParameters []string) (*lib.Job, error) {
	startJob := lib.StartJob{
		ScmReference: scmReference,
		Parameters:   envParameters,
	}
	createdJob := &lib.Job{}

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/job", url.QueryEscape(projectID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &startJob,
		Results:    &createdJob,
		OkCodes:    []int{202},
		SetHeaders: c.config.authenticateRequest,
	})

	return createdJob, err
}

func (c *Project) EncryptData(projectID, toEncryptString string) (string, error) {
	toEncryptData := &lib.StringValue{
		Value: toEncryptString,
	}

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("project/%s/crypto", url.QueryEscape(projectID)))
	if err != nil {
		return "", err
	}

	encryptedData := &lib.StringValue{}

	err = perigee.Put(requestURL, perigee.Options{
		ReqBody:    &toEncryptData,
		Results:    &encryptedData,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return encryptedData.Value, err
}

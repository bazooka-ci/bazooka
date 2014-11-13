package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	lib "github.com/bazooka-ci/bazooka-lib"
)

type Client struct {
	URL string
}

func NewClient(endpoint string) (*Client, error) {
	return &Client{
		URL: endpoint,
	}, nil
}

func (c *Client) ListProjects() ([]lib.Project, error) {
	resp, err := http.Get(fmt.Sprintf("%s/project/", c.URL))
	if err != nil {
		return nil, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return nil, err
	}
	body, err := body(resp)
	if err != nil {
		return nil, err
	}
	var projects []lib.Project
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *Client) CreateProject(name, scm, scmUri string) (*lib.Project, error) {
	project := lib.Project{
		Name:    name,
		ScmType: scm,
		ScmURI:  scmUri,
	}
	p, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/project/", c.URL), "application/json", bytes.NewBuffer(p))
	if err != nil {
		return nil, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return nil, err
	}
	body, err := body(resp)
	if err != nil {
		return nil, err
	}
	createdProject := &lib.Project{}
	err = json.Unmarshal(body, createdProject)
	if err != nil {
		return nil, err
	}
	return createdProject, nil
}

func (c *Client) StartJob(projectID, scmReference string) (*lib.Job, error) {
	startJob := lib.StartJob{
		ScmReference: scmReference,
	}
	p, err := json.Marshal(startJob)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/project/%s/job", c.URL, projectID), "application/json", bytes.NewBuffer(p))
	if err != nil {
		return nil, err
	}
	err = c.checkResponse(resp)
	if err != nil {
		return nil, err
	}
	body, err := body(resp)
	if err != nil {
		return nil, err
	}
	createJob := &lib.Job{}
	err = json.Unmarshal(body, createJob)
	if err != nil {
		return nil, err
	}
	return createJob, nil
}

func (c *Client) checkResponse(s *http.Response) error {
	switch {
	case s.StatusCode == http.StatusNotFound:
		return fmt.Errorf("Object not found")
	case s.StatusCode >= 400 && s.StatusCode < 500:
		b, err := body(s)
		if err != nil {
			return fmt.Errorf("Error (%v)", s.Status)
		}
		return fmt.Errorf("Internal error:\n%v", string(b))
	case s.StatusCode == 500:
		b, err := body(s)
		if err != nil {
			return fmt.Errorf("Internal error (%v)", s.Status)
		}
		return fmt.Errorf("Internal error:\n%v", string(b))
	case s.StatusCode >= 501 && s.StatusCode < 505:
		return fmt.Errorf("Service maintainance (%v)", s.Status)
	default:
		return nil
	}
}

func body(s *http.Response) ([]byte, error) {
	return ioutil.ReadAll(s.Body)
}

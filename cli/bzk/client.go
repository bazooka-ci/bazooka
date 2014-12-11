package main

import (
	"fmt"
	"net/http"
	"net/url"

	lib "github.com/haklop/bazooka/commons"
	"github.com/racker/perigee"
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

func (c *Client) ListVariants(jobID string) ([]lib.Variant, error) {
	var v []lib.Variant

	requestURL, err := c.getRequestURL(fmt.Sprintf("job/%s/variant", url.QueryEscape(jobID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &v,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return v, err
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

func (c *Client) StartJob(projectID, scmReference string) (*lib.Job, error) {
	startJob := lib.StartJob{
		ScmReference: scmReference,
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

func (c *Client) JobLog(jobID string) ([]lib.LogEntry, error) {
	var log []lib.LogEntry

	requestURL, err := c.getRequestURL(fmt.Sprintf("job/%s/log", url.QueryEscape(jobID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &log,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return log, err
}

func (c *Client) VariantLog(variantID string) ([]lib.LogEntry, error) {
	var log []lib.LogEntry

	requestURL, err := c.getRequestURL(fmt.Sprintf("variant/%s/log", url.QueryEscape(variantID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &log,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})
	return log, err
}

func (c *Client) ListImages() ([]*lib.Image, error) {

	var images []*lib.Image

	requestURL, err := c.getRequestURL("image")
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &images,
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return images, err
}

func (c *Client) SetImage(name, image string) error {

	requestURL, err := c.getRequestURL(fmt.Sprintf("image/%s", url.QueryEscape(name)))
	if err != nil {
		return err
	}

	err = perigee.Put(requestURL, perigee.Options{
		ReqBody: &struct {
			Image string `json:"image"`
		}{image},
		OkCodes:    []int{200},
		SetHeaders: c.authenticateRequest,
	})

	return err
}

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

func (c *Client) authenticateRequest(r *http.Request) error {
	authConfig, err := loadConfig()
	if err != nil {
		return err
	}

	if len(authConfig.Username) > 0 {
		r.SetBasicAuth(authConfig.Username, authConfig.Password)
	}
	return nil
}

func (c *Client) getRequestURL(path string) (string, error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return "", fmt.Errorf("Bazooka URL %s has an incorrect format: %v", c.URL, err)
	}
	u.Path = u.Path + "/" + path
	return u.String(), nil
}

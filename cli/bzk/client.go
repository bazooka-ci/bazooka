package main

import (
	"fmt"

	lib "github.com/haklop/bazooka/commons"
	"github.com/racker/perigee"
)

type Client struct {
	URL      string
	Username string
	Password string
}

func NewClient(endpoint string) (*Client, error) {
	return &Client{
		URL: endpoint,
	}, nil
}

func (c *Client) ListProjects() ([]lib.Project, error) {
	var p []lib.Project

	ep := fmt.Sprintf("%s/project", c.URL)
	err := perigee.Get(ep, perigee.Options{
		Results: &p,
		OkCodes: []int{200},
	})

	return p, err
}

func (c *Client) ListJobs(projectID string) ([]lib.Job, error) {
	var j []lib.Job

	ep := fmt.Sprintf("%s/project/%s/job", c.URL, projectID)
	err := perigee.Get(ep, perigee.Options{
		Results: &j,
		OkCodes: []int{200},
	})

	return j, err
}

func (c *Client) ListVariants(jobID string) ([]lib.Variant, error) {
	var v []lib.Variant
	ep := fmt.Sprintf("%s/job/%s/variant", c.URL, jobID)
	err := perigee.Get(ep, perigee.Options{
		Results: &v,
		OkCodes: []int{200},
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

	ep := fmt.Sprintf("%s/project", c.URL)
	err := perigee.Post(ep, perigee.Options{
		ReqBody: &project,
		Results: &createdProject,
		OkCodes: []int{201},
	})

	return createdProject, err
}

func (c *Client) StartJob(projectID, scmReference string) (*lib.Job, error) {
	startJob := lib.StartJob{
		ScmReference: scmReference,
	}
	createdJob := &lib.Job{}

	ep := fmt.Sprintf("%s/project/%s/job", c.URL, projectID)
	err := perigee.Post(ep, perigee.Options{
		ReqBody: &startJob,
		Results: &createdJob,
		OkCodes: []int{202},
	})

	return createdJob, err
}

func (c *Client) JobLog(jobID string) ([]lib.LogEntry, error) {
	var log []lib.LogEntry

	ep := fmt.Sprintf("%s/job/%s/log", c.URL, jobID)
	err := perigee.Get(ep, perigee.Options{
		Results: &log,
		OkCodes: []int{200},
	})

	return log, err
}

func (c *Client) VariantLog(variantID string) ([]lib.LogEntry, error) {
	var log []lib.LogEntry

	ep := fmt.Sprintf("%s/variant/%v/log", c.URL, variantID)
	err := perigee.Get(ep, perigee.Options{
		Results: &log,
		OkCodes: []int{200},
	})
	return log, err
}

func (c *Client) ListImages() ([]*lib.Image, error) {

	var images []*lib.Image

	ep := fmt.Sprintf(fmt.Sprintf("%s/image", c.URL))
	err := perigee.Get(ep, perigee.Options{
		Results: &images,
		OkCodes: []int{200},
	})

	return images, err
}

func (c *Client) SetImage(name, image string) error {

	ep := fmt.Sprintf("%s/image/%s", c.URL, name)
	err := perigee.Put(ep, perigee.Options{
		ReqBody: &struct {
			Image string `json:"image"`
		}{image},
		OkCodes: []int{200},
	})

	return err
}

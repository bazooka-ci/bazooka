package client

import (
	"fmt"
	"net/url"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

type Job struct {
	config *Config
}

func (c *Job) List() ([]lib.Job, error) {
	var j []lib.Job

	requestURL, err := c.config.getRequestURL("job")
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

func (c *Job) Variants(jobID string) ([]lib.Variant, error) {
	var v []lib.Variant

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("job/%s/variant", url.QueryEscape(jobID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &v,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return v, err
}

func (c *Job) Log(jobID string) ([]lib.LogEntry, error) {
	var log []lib.LogEntry

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("job/%s/log", url.QueryEscape(jobID)))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &log,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return log, err
}

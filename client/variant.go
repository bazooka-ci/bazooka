package client

import (
	"fmt"
	"net/url"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

type Variant struct {
	config *Config
}

func (c *Variant) Log(variantID string) ([]lib.LogEntry, error) {
	var log []lib.LogEntry

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("variant/%s/log", url.QueryEscape(variantID)))
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

func (c *Variant) StreamLog(variantID string) (chan lib.LogEntry, error) {
	requestURL, err := c.config.getRequestURL(fmt.Sprintf("variant/%s/log", url.QueryEscape(variantID)), "follow=true")
	if err != nil {
		return nil, err
	}

	response, err := perigee.Request("GET", requestURL, perigee.Options{
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})
	if err != nil {
		return nil, err
	}

	return streamLog(response.HttpResponse), nil
}

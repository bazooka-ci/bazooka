package client

import (
	"fmt"
	"net/url"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

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

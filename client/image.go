package client

import (
	"fmt"
	"net/url"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

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

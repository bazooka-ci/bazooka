package client

import (
	"fmt"
	"net/url"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

type Image struct {
	config *Config
}

func (c *Image) List() ([]*lib.Image, error) {

	var images []*lib.Image

	requestURL, err := c.config.getRequestURL("image")
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &images,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return images, err
}

func (c *Image) Get(name string) (*lib.Image, error) {
	var image lib.Image

	requestURL, err := c.config.getRequestURL(fmt.Sprintf("image/%s", name))
	if err != nil {
		return nil, err
	}

	err = perigee.Get(requestURL, perigee.Options{
		Results:    &image,
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return &image, err
}

func (c *Image) Set(name, image string) error {
	requestURL, err := c.config.getRequestURL(fmt.Sprintf("image/%s", url.QueryEscape(name)))
	if err != nil {
		return err
	}

	err = perigee.Put(requestURL, perigee.Options{
		ReqBody: &struct {
			Image string `json:"image"`
		}{image},
		OkCodes:    []int{200},
		SetHeaders: c.config.authenticateRequest,
	})

	return err
}

package client

import (
	"fmt"
	"net/http"
	"net/url"
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Client struct {
	config Config
}

func New(config Config) (*Client, error) {
	return &Client{
		config,
	}, nil
}

func (c *Client) authenticateRequest(r *http.Request) error {
	if len(c.config.Username) > 0 {
		r.SetBasicAuth(c.config.Username, c.config.Password)
	}
	return nil
}

func (c *Client) getRequestURL(path string) (string, error) {
	u, err := url.Parse(c.config.URL)
	if err != nil {
		return "", fmt.Errorf("Bazooka URL %s has an incorrect format: %v", c.config.URL, err)
	}
	u.Path = u.Path + "/" + path
	return u.String(), nil
}

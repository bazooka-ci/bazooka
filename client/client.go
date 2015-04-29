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
	Project *Project
	Job     *Job
	Variant *Variant
	Image   *Image
	User    *User
}

func New(config *Config) (*Client, error) {
	return &Client{
		Project: &Project{
			config: config,
			Key:    &ProjectKey{config},
			Config: &ProjectConfig{config},
		},
		Job:     &Job{config},
		Variant: &Variant{config},
		Image:   &Image{config},
		User:    &User{config},
	}, nil
}

func (c *Config) authenticateRequest(r *http.Request) error {
	if len(c.Username) > 0 {
		r.SetBasicAuth(c.Username, c.Password)
	}
	return nil
}

func (c *Config) getRequestURL(path string) (string, error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return "", fmt.Errorf("Bazooka URL %s has an incorrect format: %v", c.URL, err)
	}
	u.Path = u.Path + "/" + path
	return u.String(), nil
}

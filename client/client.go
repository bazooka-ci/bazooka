package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/racker/perigee"
)

type Config struct {
	URL      string
	Username string
	Password string
}

type Client struct {
	Project  *Project
	Job      *Job
	Variant  *Variant
	Image    *Image
	User     *User
	Internal *Internal
}

func New(config *Config) (*Client, error) {
	return &Client{
		Project: &Project{
			config: config,
			Key:    &ProjectKey{config},
			Config: &ProjectConfig{config},
		},
		Job:      &Job{config},
		Variant:  &Variant{config},
		Image:    &Image{config},
		User:     &User{config},
		Internal: &Internal{config},
	}, nil
}

func IsNotFound(err error) bool {
	switch err := err.(type) {
	case *perigee.UnexpectedResponseCodeError:
		return err.Actual == http.StatusNotFound
	default:
		return false
	}
}

func (c *Config) authenticateRequest(r *http.Request) error {
	if len(c.Username) > 0 {
		r.SetBasicAuth(c.Username, c.Password)
	}
	return nil
}

func (c *Config) getRequestURL(path string, query ...string) (string, error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return "", fmt.Errorf("Bazooka URL %s has an incorrect format: %v", c.URL, err)
	}

	u.Path = u.Path + "/" + path

	vals := url.Values{}
	seen := map[string]struct{}{}

	for _, q := range query {
		kv := strings.SplitN(q, "=", 2)
		if _, add := seen[kv[0]]; add {
			vals.Add(kv[0], kv[1])
		} else {
			vals.Set(kv[0], kv[1])
			seen[kv[0]] = struct{}{}
		}
	}

	u.RawQuery = vals.Encode()

	return u.String(), nil
}

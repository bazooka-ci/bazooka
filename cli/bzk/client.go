package main

import "github.com/bazooka-ci/bazooka/client"

func NewClient() (*client.Client, error) {
	cliConfig, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return client.New(&client.Config{
		URL:      *bzkApiUrl,
		Username: cliConfig.Username,
		Password: cliConfig.Password,
	})
}

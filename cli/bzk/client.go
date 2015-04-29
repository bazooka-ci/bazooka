package main

import (
	"fmt"
	"log"

	"github.com/bazooka-ci/bazooka/client"
)

func NewClient() (*client.Client, error) {
	cliConfig, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return client.New(client.Config{
		URL:      checkServerURI(*bzkUri, cliConfig),
		Username: cliConfig.Username,
		Password: cliConfig.Password,
	})
}

func checkServerURI(endpoint string, cliConfig *Config) string {
	if len(endpoint) == 0 {
		if len(cliConfig.ServerURI) == 0 {
			endpoint = interactiveInput("Bazooka Server URI")
			cliConfig.ServerURI = endpoint

			if err := saveConfig(cliConfig); err != nil {
				log.Fatal(fmt.Errorf("Unable to save Bazooka config, reason is: %v\n", err))
			}
		}
		return cliConfig.ServerURI
	}
	return endpoint
}

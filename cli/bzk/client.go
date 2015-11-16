package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/bazooka-ci/bazooka/client"
)

func NewClient() (*client.Client, error) {
	cliConfig, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return client.New(&client.Config{
		URL:      checkServerURI(*bzkApiUri, cliConfig),
		Username: cliConfig.Username,
		Password: cliConfig.Password,
	})
}

func checkServerURI(endpoint string, cliConfig *Config) string {
	if len(endpoint) == 0 {
		if len(cliConfig.ApiURI) == 0 {
			var defaultURI = "http://localhost:3000"
			if runtime.GOOS == "darwin" {
				defaultURI = "http://192.168.59.103:3000"
			}
			endpoint = interactiveInput("Bazooka Server URI", defaultURI)
			cliConfig.ApiURI = endpoint

			if err := saveConfig(cliConfig); err != nil {
				log.Fatal(fmt.Errorf("Unable to save Bazooka config, reason is: %v\n", err))
			}
		}
		return cliConfig.ApiURI
	}
	return endpoint
}

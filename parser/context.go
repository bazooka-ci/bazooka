package main

import (
	"encoding/json"
	"os"

	"fmt"

	"log"

	"github.com/bazooka-ci/bazooka/client"
	lib "github.com/bazooka-ci/bazooka/commons"
)

const (
	BazookaEnvApiUrl        = "BZK_API_URL"
	BazookaEnvSyslogUrl     = "BZK_SYSLOG_URL"
	BazookaEnvHome          = "BZK_HOME"
	BazookaEnvSrc           = "BZK_SRC"
	BazookaEnvCryptoKeyfile = "BZK_CRYPTO_KEYFILE"
	BazookaEnvDockerSock    = "BZK_DOCKERSOCK"
	BazookaEnvProjectID     = "BZK_PROJECT_ID"
	BazookaEnvJobID         = "BZK_JOB_ID"
	BazookaEnvJobParameters = "BZK_JOB_PARAMETERS"
)

type context struct {
	client        *client.Client
	apiUrl        string
	syslogUrl     string
	network       string
	projectID     string
	jobID         string
	jobParameters string
	paths         paths
}

type paths struct {
	source         path
	output         path
	meta           path
	cryptoKey      path
	dockerSock     path
	dockerEndpoint path
}

type path struct {
	container string
	host      string
}

func initContext() *context {
	apiUrl := os.Getenv(BazookaEnvApiUrl)
	// Configure Client
	client, err := client.New(&client.Config{
		URL: apiUrl,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &context{
		client:        client,
		apiUrl:        apiUrl,
		syslogUrl:     os.Getenv(BazookaEnvSyslogUrl),
		projectID:     os.Getenv(BazookaEnvProjectID),
		jobID:         os.Getenv(BazookaEnvJobID),
		jobParameters: os.Getenv(BazookaEnvJobParameters),
		paths: paths{
			source:         path{"/bazooka", os.Getenv(BazookaEnvSrc)},
			output:         path{"/bazooka-output", os.Getenv(BazookaEnvHome) + "/work"},
			meta:           path{"/meta", os.Getenv(BazookaEnvHome) + "/meta"},
			cryptoKey:      path{"/bazooka-cryptokey", os.Getenv(BazookaEnvCryptoKeyfile)},
			dockerSock:     path{"/var/run/docker.sock", ""},
			dockerEndpoint: path{"unix:///var/run/docker.sock", ""},
		},
	}
}

func (c *context) loggerConfig(image string) map[string]string {
	return map[string]string{
		"syslog-address": c.syslogUrl,
		"syslog-tag":     fmt.Sprintf("image=%s;project=%s;job=%s", image, c.projectID, c.jobID),
	}
}

func (c *context) unmarshalJobParameters() ([]lib.BzkString, error) {
	var res []lib.BzkString
	err := json.Unmarshal([]byte(c.jobParameters), &res)
	return res, err
}

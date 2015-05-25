package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bazooka-ci/bazooka/client"
	lib "github.com/bazooka-ci/bazooka/commons"
)

const (
	BazookaEnvHome          = "BZK_HOME"
	BazookaEnvSrc           = "BZK_SRC"
	BazookaEnvCryptoKeyfile = "BZK_CRYPTO_KEYFILE"
	BazookaEnvDockerSock    = "BZK_DOCKERSOCK"
	BazookaEnvProjectID     = "BZK_PROJECT_ID"
	BazookaEnvJobID         = "BZK_JOB_ID"
	BazookaEnvJobParameters = "BZK_JOB_PARAMETERS"

	BazookaEnvServerAddr = "SERVER_PORT_3000_TCP_ADDR"
	BazookaEnvServerPort = "SERVER_PORT_3000_TCP_PORT"

	BazookaEnvLogServerAddr = "SERVER_PORT_3001_TCP_ADDR"
	BazookaEnvLogServerPort = "SERVER_PORT_3001_TCP_PORT"
)

type context struct {
	client        *client.Client
	projectID     string
	jobID         string
	jobParameters string
	syslogAddress string
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
	// Configure Client
	client, err := client.New(&client.Config{
		URL: fmt.Sprintf("http://%s:%s", os.Getenv(BazookaEnvServerAddr), os.Getenv(BazookaEnvServerPort)),
	})
	if err != nil {
		log.Fatal(err)
	}

	return &context{
		client:        client,
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
		"syslog-address": fmt.Sprintf("tcp://%s:%s", os.Getenv(BazookaEnvLogServerAddr), os.Getenv(BazookaEnvLogServerPort)),
		"syslog-tag":     fmt.Sprintf("image=%s;project=%s;job=%s", image, c.projectID, c.jobID),
	}
}

func (c *context) unmarshalJobParameters() ([]lib.BzkString, error) {
	var res []lib.BzkString
	err := json.Unmarshal([]byte(c.jobParameters), &res)
	return res, err
}

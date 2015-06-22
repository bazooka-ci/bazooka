package main

import (
	"encoding/json"
	"os"

	lib "github.com/bazooka-ci/bazooka/commons"
)

const (
	BazookaEnvHome          = "BZK_HOME"
	BazookaEnvSrc           = "BZK_SRC"
	BazookaEnvCryptoKeyfile = "BZK_CRYPTO_KEYFILE"
	BazookaEnvDockerSock    = "BZK_DOCKERSOCK"
	BazookaEnvProjectID     = "BZK_PROJECT_ID"
	BazookaEnvJobID         = "BZK_JOB_ID"
)

type context struct {
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
	return &context{
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

func (c *context) unmarshalJobParameters() ([]lib.BzkString, error) {
	var res []lib.BzkString
	err := json.Unmarshal([]byte(c.jobParameters), &res)
	return res, err
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bazooka-ci/bazooka/client"
)

const (
	BazookaEnvServerName    = "BZK_SERVER_NAME"
	BazookaEnvHome          = "BZK_HOME"
	BazookaEnvSrc           = "BZK_SRC"
	BazookaEnvSCMKeyfile    = "BZK_SCM_KEYFILE"
	BazookaEnvCryptoKeyfile = "BZK_CRYPTO_KEYFILE"
	BazookaEnvDockerSock    = "BZK_DOCKERSOCK"
	BazookaEnvSCM           = "BZK_SCM"
	BazookaEnvSCMUrl        = "BZK_SCM_URL"
	BazookaEnvSCMReference  = "BZK_SCM_REFERENCE"
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
	serverName    string
	scm           string
	scmUrl        string
	scmReference  string
	projectID     string
	jobID         string
	jobParameters string
	reuseScm      bool
	paths         paths
}

type paths struct {
	base           path
	source         path
	work           path
	meta           path
	artifacts      path
	scmKey         path
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
		serverName:    os.Getenv(BazookaEnvServerName),
		scm:           os.Getenv(BazookaEnvSCM),
		scmUrl:        os.Getenv(BazookaEnvSCMUrl),
		scmReference:  os.Getenv(BazookaEnvSCMReference),
		projectID:     os.Getenv(BazookaEnvProjectID),
		jobID:         os.Getenv(BazookaEnvJobID),
		jobParameters: os.Getenv(BazookaEnvJobParameters),
		reuseScm:      os.Getenv("BZK_REUSE_SCM_CHECKOUT") != "",
		paths: paths{
			base:           path{"/bazooka", os.Getenv(BazookaEnvHome)},
			source:         path{"/bazooka/source", os.Getenv(BazookaEnvSrc)},
			work:           path{"/bazooka/work", os.Getenv(BazookaEnvHome) + "/work"},
			meta:           path{"/bazooka/meta", os.Getenv(BazookaEnvHome) + "/meta"},
			artifacts:      path{"/bazooka/artifacts", os.Getenv(BazookaEnvHome) + "/artifacts"},
			scmKey:         path{"/bazooka/key", os.Getenv(BazookaEnvSCMKeyfile)},
			cryptoKey:      path{"/bazooka/crypto-key", os.Getenv(BazookaEnvCryptoKeyfile)},
			dockerSock:     path{"/var/run/docker.sock", os.Getenv(BazookaEnvDockerSock)},
			dockerEndpoint: path{"unix:///var/run/docker.sock", "unix://" + os.Getenv(BazookaEnvDockerSock)},
		},
	}
}

func (c *context) loggerConfig(image string, variantID string) map[string]string {
	tag := fmt.Sprintf("image=%s;project=%s;job=%s", image, c.projectID, c.jobID)
	if len(variantID) > 0 {
		tag += ";variant=" + variantID
	}
	return map[string]string{
		"syslog-address": fmt.Sprintf("tcp://%s:%s", os.Getenv(BazookaEnvLogServerAddr), os.Getenv(BazookaEnvLogServerPort)),
		"syslog-tag":     tag,
	}
}

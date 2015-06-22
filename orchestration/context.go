package main

import (
	"os"

	"github.com/bazooka-ci/bazooka/commons/mongo"
)

const (
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
)

type context struct {
	connector     *mongo.MongoConnector
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

	return &context{
		connector:     mongo.NewConnector(),
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

func (c *context) cleanup() {
	c.connector.Close()
}

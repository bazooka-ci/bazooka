package main

import (
	"os"
)

const (
	BazookaEnvHome         = "BZK_HOME"
	BazookaEnvSCMKeyfile   = "BZK_SCM_KEYFILE"
	BazookaEnvDockerSock   = "BZK_DOCKERSOCK"
	BazookaEnvSCM          = "BZK_SCM"
	BazookaEnvSCMUrl       = "BZK_SCM_URL"
	BazookaEnvSCMReference = "BZK_SCM_REFERENCE"
	BazookaEnvProjectID    = "BZK_PROJECT_ID"
	BazookaEnvJobID        = "BZK_JOB_ID"
)

type stdPaths struct {
	base           string
	source         string
	work           string
	meta           string
	artifacts      string
	key            string
	cryptoKey      string
	dockerSock     string
	dockerEndpoint string
}

type bzkPaths struct {
	container stdPaths
	host      stdPaths
}

var paths = bzkPaths{
	stdPaths{
		"/bazooka",
		"/bazooka/source",
		"/bazooka/work",
		"/bazooka/meta",
		"/bazooka/artifacts",
		"/bazooka/key",
		"/bazooka-cryptokey",
		"/var/run/docker.sock",
		"unix:///var/run/docker.sock",
	},
	stdPaths{
		os.Getenv(BazookaEnvHome),
		os.Getenv(BazookaEnvHome) + "/source",
		os.Getenv(BazookaEnvHome) + "/work",
		os.Getenv(BazookaEnvHome) + "/meta",
		os.Getenv(BazookaEnvHome) + "/artifacts",
		os.Getenv(BazookaEnvSCMKeyfile),
		os.Getenv(BazookaEnvHome) + "/crypto-key",
		os.Getenv(BazookaEnvDockerSock),
		"unix://" + os.Getenv(BazookaEnvDockerSock),
	},
}

package main

import (
	"os"
)

const (
	BazookaEnvHome       = "BZK_HOME"
	BazookaEnvSCMKeyfile = "BZK_SCM_KEYFILE"
	BazookaEnvDockerSock = "BZK_DOCKERSOCK"
)

type stdPaths struct {
	source         string
	meta           string
	output         string
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
		"/meta",
		"/bazooka-output",
		"/docker.sock",
		"unix:///docker.sock",
	},
	stdPaths{
		os.Getenv(BazookaEnvHome) + "/source",
		os.Getenv(BazookaEnvHome) + "/meta",
		os.Getenv(BazookaEnvHome) + "/work",
		"",
		"",
	},
}

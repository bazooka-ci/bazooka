package main

import (
	"os"
)

const (
	BazookaEnvHome       = "BZK_HOME"
	BazookaEnvSrc        = "BZK_SRC"
	BazookaEnvSCMKeyfile = "BZK_SCM_KEYFILE"
	BazookaEnvDockerSock = "BZK_DOCKERSOCK"
)

type stdPaths struct {
	source         string
	meta           string
	output         string
	dockerSock     string
	dockerEndpoint string
	cryptoKey      string
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
		"/bazooka-cryptokey",
	},
	stdPaths{
		os.Getenv(BazookaEnvSrc),
		os.Getenv(BazookaEnvHome) + "/meta",
		os.Getenv(BazookaEnvHome) + "/work",
		"",
		"",
		os.Getenv(BazookaEnvHome) + "/crypto-key",
	},
}

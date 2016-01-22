package main

import (
	"log"
	"os"
	"time"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/commons/mongo"
)

const (
	BazookaEnvApiUrl     = "BZK_API_URL"
	BazookaEnvSyslogUrl  = "BZK_SYSLOG_URL"
	BazookaEnvDbUrl      = "BZK_DB_URL"
	BazookaEnvNetwork    = "BZK_NETWORK"
	BazookaEnvSCMKeyfile = "BZK_SCM_KEYFILE"
	BazookaEnvHome       = "BZK_HOME"
	BazookaEnvDockerSock = "BZK_DOCKERSOCK"

	DockerSock     = "/var/run/docker.sock"
	DockerEndpoint = "unix://" + DockerSock
	BazookaHome    = "/bazooka"
)

type context struct {
	apiUrl    string
	syslogUrl string
	dbUrl     string
	network   string
	connector *mongo.MongoConnector
	paths     paths
}

type paths struct {
	home           path
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
	c := &context{
		apiUrl:    os.Getenv(BazookaEnvApiUrl),
		syslogUrl: os.Getenv(BazookaEnvSyslogUrl),
		dbUrl:     os.Getenv(BazookaEnvDbUrl),
		network:   os.Getenv(BazookaEnvNetwork),
		paths: paths{
			home:           path{BazookaHome, os.Getenv(BazookaEnvHome)},
			scmKey:         path{"", os.Getenv(BazookaEnvSCMKeyfile)},
			dockerSock:     path{DockerSock, os.Getenv(BazookaEnvDockerSock)},
			dockerEndpoint: path{DockerEndpoint, "unix://" + os.Getenv(BazookaEnvDockerSock)},
		},
	}

	if err := lib.WaitForTcpConnection(c.dbUrl, 100*time.Millisecond, 5*time.Second); err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}
	c.connector = mongo.NewConnector(c.dbUrl)
	return c
}

func (c *context) cleanup() {
	c.connector.Close()
}

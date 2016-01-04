package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	lib "github.com/bazooka-ci/bazooka/commons"

	"sync"

	"github.com/bazooka-ci/bazooka/client"
	docker "github.com/bywan/go-dockercommand"
	"github.com/iwanbk/gobeanstalk"
)

const (
	BazookaEnvServerApi    = "BZK_API_URL"
	BazookaEnvServerSyslog = "BZK_SYSLOG_URL"
	BazookaEnvQueue        = "BZK_QUEUE_URL"
	BazookaEnvHome         = "BZK_HOME"
	BazookaEnvSCMKeyfile   = "BZK_SCM_KEYFILE"
	BazookaEnvDockerSock   = "BZK_DOCKERSOCK"
	BazookaEnvNetwork      = "BZK_NETWORK"

	DockerSock = "/var/run/docker.sock"
)

type context struct {
	client       *client.Client
	queue        *gobeanstalk.Conn
	docker       *docker.Docker
	busy         sync.WaitGroup
	serverApi    string
	serverSyslog string
	network      string
	slots        int
	paths        paths
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
	// Configure Client
	client, err := client.New(&client.Config{
		URL: fmt.Sprintf(os.Getenv(BazookaEnvServerApi)),
	})
	if err != nil {
		log.Fatal(err)
	}

	queueUrl := os.Getenv(BazookaEnvQueue)
	if err := lib.WaitForTcpConnection(queueUrl, 500*time.Millisecond, 30*time.Second); err != nil {
		log.Fatalf("Cannot connect to the queue @ %s: %v", queueUrl, err)
	}
	queue, err := gobeanstalk.Dial(queueUrl)
	if err != nil {
		log.Fatal(err)
	}

	dockerClient, err := docker.NewDocker("unix://" + DockerSock)
	if err != nil {
		log.Fatal(err)
	}

	return &context{
		client:       client,
		queue:        queue,
		docker:       dockerClient,
		serverApi:    os.Getenv(BazookaEnvServerApi),
		serverSyslog: os.Getenv(BazookaEnvServerSyslog),
		network:      os.Getenv(BazookaEnvNetwork),
		paths: paths{
			home:           path{"/bazooka", os.Getenv(BazookaEnvHome)},
			scmKey:         path{"", os.Getenv(BazookaEnvSCMKeyfile)},
			dockerSock:     path{DockerSock, os.Getenv(BazookaEnvDockerSock)},
			dockerEndpoint: path{"unix://" + DockerSock, "unix://" + os.Getenv(BazookaEnvDockerSock)},
		},
	}
}

func (c *context) reserveJob() (*reservedJob, error) {
	j, err := c.queue.Reserve()
	if err != nil {
		return nil, fmt.Errorf("Error while reserving a job: %v", err)
	}
	var job lib.Job
	if err := json.Unmarshal(j.Body, &job); err != nil {
		return nil, fmt.Errorf("Error while parsing a job: %v", err)
	}
	return &reservedJob{id: j.ID, job: &job}, nil
}

func (c *context) deleteJob(id uint64) error {
	if err := c.queue.Delete(id); err != nil {
		return fmt.Errorf("Error while deleting job %v: %v", id, err)
	}
	return nil
}

func (c *context) releaseJob(id uint64) error {
	if err := c.queue.Release(id, 0, 30*time.Second); err != nil {
		return fmt.Errorf("Error while releasing job %v: %v", id, err)
	}
	return nil
}

func (c *context) touchJob(id uint64) error {
	if err := c.queue.Touch(id); err != nil {
		return fmt.Errorf("Error while touching job %v: %v", id, err)
	}
	return nil
}

const (
	buildFolderPattern        = "%s/build/%s/%s"     // $bzk_home/build/$projectId/$buildId
	sharedSourceFolderPattern = "%s/build/%s/source" // $bzk_home/build/$projectId/source
	logFolderPattern          = "%s/build/%s/%s/log" // $bzk_home/build/$projectId/$buildId/log
)

func (c *context) buildFolder(job *lib.Job) path {
	return path{
		host:      fmt.Sprintf(buildFolderPattern, c.paths.home.host, job.ProjectID, job.ID),
		container: fmt.Sprintf(buildFolderPattern, c.paths.home.container, job.ProjectID, job.ID),
	}
}

package main

import (
	"log"
	"os"
	"time"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/server/db"
	"github.com/iwanbk/gobeanstalk"
)

const (
	BazookaEnvQueueUrl = "BZK_QUEUE_URL"
	BazookaEnvDbUrl    = "BZK_DB_URL"
)

type context struct {
	queue     *gobeanstalk.Conn
	connector *db.MongoConnector
}

func initContext() *context {
	queueUrl := os.Getenv(BazookaEnvQueueUrl)
	if err := lib.WaitForTcpConnection(queueUrl, 500*time.Millisecond, 30*time.Second); err != nil {
		log.Fatalf("Cannot connect to the queue (%s): %v", queueUrl, err)
	}
	queue, err := gobeanstalk.Dial(queueUrl)
	if err != nil {
		log.Fatalf("Cannot connect to the queue: %v", err)
	}

	dbUrl := os.Getenv(BazookaEnvDbUrl)
	if err := lib.WaitForTcpConnection(dbUrl, 500*time.Millisecond, 30*time.Second); err != nil {
		log.Fatalf("Cannot connect to the database (%s): %v", dbUrl, err)
	}

	c := &context{
		queue:     queue,
		connector: db.NewConnector(dbUrl),
	}

	return c
}

func (c *context) cleanup() {
	c.connector.Close()
}

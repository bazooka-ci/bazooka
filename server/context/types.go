package context

import (
	"github.com/bazooka-ci/bazooka-lib/mongo"
)

type ErrorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

type Context struct {
	Connector      *mongo.MongoConnector
	DockerEndpoint string
	Env            map[string]string
}

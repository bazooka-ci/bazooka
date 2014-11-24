package context

import (
	"github.com/haklop/bazooka/commons/mongo"
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

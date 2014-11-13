package context

import (
	"gopkg.in/mgo.v2"
)

type ErrorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

type Context struct {
	Database       *mgo.Database
	DockerEndpoint string
	Env            map[string]string
}

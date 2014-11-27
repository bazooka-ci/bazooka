package main

import (
	"encoding/json"
	"net/http"

	"github.com/haklop/bazooka/commons/mongo"
)

const (
	BazookaEnvSCMKeyfile = "BZK_SCM_KEYFILE"
	BazookaEnvHome       = "BZK_HOME"
	BazookaEnvDockerSock = "BZK_DOCKERSOCK"
	BazookaEnvMongoAddr  = "MONGO_PORT_27017_TCP_ADDR"
	BazookaEnvMongoPort  = "MONGO_PORT_27017_TCP_PORT"

	DockerSock     = "/var/run/docker.sock"
	DockerEndpoint = "unix://" + DockerSock
	BazookaHome    = "/bazooka"
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

func WriteError(err error, res http.ResponseWriter, encoder *json.Encoder) {
	res.WriteHeader(500)
	encoder.Encode(&ErrorResponse{
		Code:    500,
		Message: err.Error(),
	})
}

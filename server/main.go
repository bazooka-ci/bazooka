package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bazooka-ci/bazooka-lib/mongo"
	"github.com/gorilla/mux"
	"github.com/haklop/bazooka/server/context"
	"github.com/haklop/bazooka/server/fetcher"
	"github.com/haklop/bazooka/server/project"
)

func main() {
	// Configure Bazooka
	env := map[string]string{
		context.BazookaEnvSCMKeyfile: os.Getenv(context.BazookaEnvSCMKeyfile),
		context.BazookaEnvHome:       os.Getenv(context.BazookaEnvHome),
		context.BazookaEnvDockerSock: os.Getenv(context.BazookaEnvDockerSock),
		context.BazookaEnvMongoAddr:  os.Getenv(context.BazookaEnvMongoAddr),
		context.BazookaEnvMongoPort:  os.Getenv(context.BazookaEnvMongoPort),
	}

	if len(env[context.BazookaEnvHome]) == 0 {
		env[context.BazookaEnvHome] = "/bazooka"
	}

	// Enable bazooka-server to be execute without running its own container
	var serverEndpoint string
	if len(os.Getenv("DOCKER_HOST")) != 0 {
		serverEndpoint = os.Getenv("DOCKER_HOST")
	} else {
		serverEndpoint = context.DockerEndpoint
	}

	// Configure Mongo
	connector := mongo.NewConnector()
	defer connector.Close()

	serverContext := context.Context{
		Connector:      connector,
		DockerEndpoint: serverEndpoint,
		Env:            env,
	}

	// Configure web server
	r := mux.NewRouter()

	projectRouter := r.PathPrefix("/project").Subrouter()
	projectHandler := project.Handlers{}
	projectHandler.SetHandlers(projectRouter, serverContext)

	fetcherRouter := r.PathPrefix("/fetcher").Subrouter()
	fetcherHandler := fetcher.Handlers{}
	fetcherHandler.SetHandlers(fetcherRouter, serverContext)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

package main

import (
	"log"
	"net/http"
	"os"

	"bitbucket.org/bywan/bazooka-api/bazooka-server/server/context"
	"bitbucket.org/bywan/bazooka-api/bazooka-server/server/fetcher"
	"bitbucket.org/bywan/bazooka-api/bazooka-server/server/project"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
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
	session, err := mgo.Dial(env[context.BazookaEnvMongoAddr] + ":" + env[context.BazookaEnvMongoPort])
	if err != nil {
		panic(err)
	}
	defer session.Close()

	database := session.DB("bazooka")

	serverContext := context.Context{
		Database:       database,
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

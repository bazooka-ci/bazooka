package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/haklop/bazooka/commons/mongo"
)

func main() {
	// Configure Bazooka
	env := map[string]string{
		BazookaEnvSCMKeyfile: os.Getenv(BazookaEnvSCMKeyfile),
		BazookaEnvHome:       os.Getenv(BazookaEnvHome),
		BazookaEnvDockerSock: os.Getenv(BazookaEnvDockerSock),
		BazookaEnvMongoAddr:  os.Getenv(BazookaEnvMongoAddr),
		BazookaEnvMongoPort:  os.Getenv(BazookaEnvMongoPort),
	}

	if len(env[BazookaEnvHome]) == 0 {
		env[BazookaEnvHome] = "/bazooka"
	}

	// Enable bazooka-server to be execute without running its own container
	var serverEndpoint string
	if len(os.Getenv("DOCKER_HOST")) != 0 {
		serverEndpoint = os.Getenv("DOCKER_HOST")
	} else {
		serverEndpoint = DockerEndpoint
	}

	// Configure Mongo
	connector := mongo.NewConnector()
	defer connector.Close()

	ctx := Context{
		Connector:      connector,
		DockerEndpoint: serverEndpoint,
		Env:            env,
	}

	// Configure web server
	r := mux.NewRouter()

	r.HandleFunc("/project", ctx.createProject).Methods("POST")
	r.HandleFunc("/project", ctx.getProjects).Methods("GET")
	r.HandleFunc("/project/{id}", ctx.getProject).Methods("GET")
	r.HandleFunc("/project/{id}/job", ctx.startBuild).Methods("POST")
	r.HandleFunc("/project/{id}/job", ctx.getJobs).Methods("GET")

	r.HandleFunc("/job/{id}", ctx.getJob).Methods("GET")
	r.HandleFunc("/job/{id}/log", ctx.getJobLog).Methods("GET")
	r.HandleFunc("/job/{id}/variant", ctx.getVariants).Methods("GET")

	r.HandleFunc("/variant/{id}", ctx.getVariant).Methods("GET")
	r.HandleFunc("/variant/{id}/log", ctx.getVariantLog).Methods("GET")

	r.HandleFunc("/fetcher", ctx.createFetcher).Methods("POST")
	r.HandleFunc("/fetcher", ctx.getFetchers).Methods("GET")
	r.HandleFunc("/fetcher/{id}", ctx.getFetcher).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

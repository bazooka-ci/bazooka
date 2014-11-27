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

	p := Context{
		Connector:      connector,
		DockerEndpoint: serverEndpoint,
		Env:            env,
	}

	// Configure web server
	r := mux.NewRouter()

	r.HandleFunc("/project", p.createProject).Methods("POST")
	r.HandleFunc("/project", p.getProjects).Methods("GET")
	r.HandleFunc("/project/{id}", p.getProject).Methods("GET")
	r.HandleFunc("/project/{id}/job", p.startBuild).Methods("POST")
	r.HandleFunc("/project/{id}/job", p.getJobs).Methods("GET")

	r.HandleFunc("/job/{job_id}", p.getJob).Methods("GET")
	r.HandleFunc("/job/{job_id}/log", p.getJobLog).Methods("GET")
	r.HandleFunc("/job/{job_id}/variant", p.getVariants).Methods("GET")

	r.HandleFunc("/variant/{variant_id}", p.getVariant).Methods("GET")
	r.HandleFunc("/variant/{variant_id}/log", p.getVariantLog).Methods("GET")

	r.HandleFunc("/fetcher", p.createFetcher).Methods("POST")
	r.HandleFunc("/fetcher", p.getFetchers).Methods("GET")
	r.HandleFunc("/fetcher/{id}", p.getFetcher).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

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

	ctx := context{
		Connector:      connector,
		DockerEndpoint: serverEndpoint,
		Env:            env,
	}

	// Configure web server
	r := mux.NewRouter()

	r.HandleFunc("/project", mkHandler(ctx.createProject)).Methods("POST")

	r.HandleFunc("/project", mkHandler(ctx.getProjects)).Methods("GET")
	r.HandleFunc("/project/{id}", mkHandler(ctx.getProject)).Methods("GET")
	r.HandleFunc("/project/{id}/job", mkHandler(ctx.startBuild)).Methods("POST")
	r.HandleFunc("/project/{id}/job", mkHandler(ctx.getJobs)).Methods("GET")

	r.HandleFunc("/job/{id}", mkHandler(ctx.getJob)).Methods("GET")
	r.HandleFunc("/job/{id}/log", mkHandler(ctx.getJobLog)).Methods("GET")
	r.HandleFunc("/job/{id}/variant", mkHandler(ctx.getVariants)).Methods("GET")

	r.HandleFunc("/variant/{id}", mkHandler(ctx.getVariant)).Methods("GET")
	r.HandleFunc("/variant/{id}/log", mkHandler(ctx.getVariantLog)).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

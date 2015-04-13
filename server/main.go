package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
	"github.com/bazooka-ci/bazooka/commons/mongo"
	"github.com/gorilla/mux"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

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

	if err := ensureDefaultImagesExist(connector); err != nil {
		log.Fatal(err)
	}

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
	r.HandleFunc("/project/{id}/job", mkHandler(ctx.startStandardJob)).Methods("POST")
	r.HandleFunc("/project/{id}/bitbucket", mkHandler(ctx.startBitbucketJob)).Methods("POST")
	r.HandleFunc("/project/{id}/github", mkHandler(ctx.startGithubJob)).Methods("POST")
	r.HandleFunc("/project/{id}/job", mkHandler(ctx.getJobs)).Methods("GET")

	r.HandleFunc("/project/{id}/config", mkHandler(ctx.getProjectConfig)).Methods("GET")
	r.HandleFunc("/project/{id}/config/{key}", ctx.setProjectConfigKey).Methods("PUT")
	r.HandleFunc("/project/{id}/config/{key}", mkHandler(ctx.unsetProjectConfigKey)).Methods("DELETE")

	r.HandleFunc("/project/{id}/key", mkHandler(ctx.addKey)).Methods("POST")
	r.HandleFunc("/project/{id}/key", mkHandler(ctx.updateKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/key", mkHandler(ctx.listKeys)).Methods("GET")

	r.HandleFunc("/job", mkHandler(ctx.getAllJobs)).Methods("GET")
	r.HandleFunc("/job/{id}", mkHandler(ctx.getJob)).Methods("GET")
	r.HandleFunc("/job/{id}/log", mkHandler(ctx.getJobLog)).Methods("GET")
	r.HandleFunc("/job/{id}/variant", mkHandler(ctx.getVariants)).Methods("GET")

	r.HandleFunc("/variant/{id}", mkHandler(ctx.getVariant)).Methods("GET")
	r.HandleFunc("/variant/{id}/log", mkHandler(ctx.getVariantLog)).Methods("GET")
	r.HandleFunc("/variant/{id}/artifacts/{path:.*}", ctx.getVariantArtifacts).Methods("GET")

	r.HandleFunc("/image", mkHandler(ctx.getImages)).Methods("GET")
	r.HandleFunc("/image/{name:.*}", mkHandler(ctx.setImage)).Methods("PUT")
	// r.HandleFunc("/image/{name}", mkHandler(ctx.unsetImage)).Methods("DELETE")

	r.HandleFunc("/user", mkHandler(ctx.getUsers)).Methods("GET")
	r.HandleFunc("/user", mkHandler(ctx.createUser)).Methods("POST")
	r.HandleFunc("/user/{id}", mkHandler(ctx.getUser)).Methods("GET")

	http.Handle("/", ctx.authenticationHandler(r))
	log.Fatal(http.ListenAndServe(":3000", nil))
}

var (
	defaultImages = map[string]string{
		"orchestration": "bazooka/orchestration",
		"parser":        "bazooka/parser",
		"scm/fetch/git": "bazooka/scm-git",
		"scm/fetch/hg":  "bazooka/scm-hg",
	}
)

func ensureDefaultImagesExist(c *mongo.MongoConnector) error {
	for name, image := range defaultImages {
		exist, err := c.HasImage(name)
		if err != nil {
			return err
		}
		if !exist {
			if err := c.SetImage(name, image); err != nil {
				return fmt.Errorf("Error while registering %s:%s: %v", name, image, err)
			}
		}

	}
	return nil
}

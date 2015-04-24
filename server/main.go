package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	lib "github.com/bazooka-ci/bazooka/commons"

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
	if err := lib.WaitForTcpConnection(env[BazookaEnvMongoAddr], env[BazookaEnvMongoPort], 100*time.Millisecond, 5*time.Second); err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}

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

	r.HandleFunc("/project", ctx.mkAuthHandler(ctx.createProject)).Methods("POST")

	r.HandleFunc("/project", ctx.mkAuthHandler(ctx.getProjects)).Methods("GET")
	r.HandleFunc("/project/{id}", ctx.mkAuthHandler(ctx.getProject)).Methods("GET")
	r.HandleFunc("/project/{id}/job", ctx.mkAuthHandler(ctx.startStandardJob)).Methods("POST")
	r.HandleFunc("/project/{id}/bitbucket", ctx.mkAuthHandler(ctx.startBitbucketJob)).Methods("POST")
	r.HandleFunc("/project/{id}/job", ctx.mkAuthHandler(ctx.getJobs)).Methods("GET")

	r.HandleFunc("/project/{id}/config", ctx.mkAuthHandler(ctx.getProjectConfig)).Methods("GET")
	r.HandleFunc("/project/{id}/config/{key}", ctx.mkAuthHandler(ctx.setProjectConfigKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/config/{key}", ctx.mkAuthHandler(ctx.unsetProjectConfigKey)).Methods("DELETE")

	r.HandleFunc("/project/{id}/key", ctx.mkAuthHandler(ctx.addKey)).Methods("POST")
	r.HandleFunc("/project/{id}/key", ctx.mkAuthHandler(ctx.updateKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/key", ctx.mkAuthHandler(ctx.listKeys)).Methods("GET")

	r.HandleFunc("/project/{id}/env", ctx.mkAuthHandler(ctx.getProjectEnv)).Methods("GET")
	r.HandleFunc("/project/{id}/env/{key}", ctx.mkAuthHandler(ctx.setProjectEnvKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/env/{key}", ctx.mkAuthHandler(ctx.unsetProjectEnvKey)).Methods("DELETE")

	r.HandleFunc("/project/{id}/crypto", ctx.mkAuthHandler(ctx.encryptData)).Methods("PUT")

	r.HandleFunc("/job", ctx.mkAuthHandler(ctx.getAllJobs)).Methods("GET")
	r.HandleFunc("/job/{id}", ctx.mkAuthHandler(ctx.getJob)).Methods("GET")
	r.HandleFunc("/job/{id}/log", ctx.mkAuthHandler(ctx.getJobLog)).Methods("GET")
	r.HandleFunc("/job/{id}/variant", ctx.mkAuthHandler(ctx.getVariants)).Methods("GET")

	r.HandleFunc("/variant/{id}", ctx.mkAuthHandler(ctx.getVariant)).Methods("GET")
	r.HandleFunc("/variant/{id}/log", ctx.mkAuthHandler(ctx.getVariantLog)).Methods("GET")
	r.HandleFunc("/variant/{id}/artifacts/{path:.*}", ctx.mkAuthHandler(ctx.getVariantArtifact)).Methods("GET")

	r.HandleFunc("/image", ctx.mkAuthHandler(ctx.getImages)).Methods("GET")
	r.HandleFunc("/image/{name:.*}", ctx.mkAuthHandler(ctx.setImage)).Methods("PUT")
	// r.HandleFunc("/image/{name}", ctx.mkAuthHandler(ctx.unsetImage)).Methods("DELETE")

	r.HandleFunc("/user", ctx.mkAuthHandler(ctx.getUsers)).Methods("GET")
	r.HandleFunc("/user", ctx.mkAuthHandler(ctx.createUser)).Methods("POST")
	r.HandleFunc("/user/{id}", ctx.mkAuthHandler(ctx.getUser)).Methods("GET")

	r.HandleFunc("/project/{id}/github", ctx.mkGithubAuthHandler(ctx.startGithubJob)).Methods("POST")

	http.Handle("/", r)

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

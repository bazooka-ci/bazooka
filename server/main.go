package main

import (
	"fmt"
	"net/http"

	"os"
	"os/signal"

	"syscall"

	log "github.com/Sirupsen/logrus"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
	"github.com/bazooka-ci/bazooka/server/db"
	"github.com/gorilla/mux"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

func main() {
	// Configure Bazooka
	context := initContext()
	defer context.cleanup()

	if err := ensureDefaultImagesExist(context.connector); err != nil {
		log.Fatal(err)
	}

	// Configure web server
	r := mux.NewRouter()

	r.HandleFunc("/project", context.mkAuthHandler(context.createProject)).Methods("POST")

	r.HandleFunc("/project", context.mkAuthHandler(context.getProjects)).Methods("GET")
	r.HandleFunc("/project/{id}", context.mkAuthHandler(context.getProject)).Methods("GET")
	r.HandleFunc("/project/{id}/job", context.mkAuthHandler(context.startStandardJob)).Methods("POST")
	r.HandleFunc("/project/{id}/job", context.mkAuthHandler(context.getJobs)).Methods("GET")

	r.HandleFunc("/project/{id}/config", context.mkAuthHandler(context.getProjectConfig)).Methods("GET")
	r.HandleFunc("/project/{id}/config/{key}", context.mkAuthHandler(context.setProjectConfigKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/config/{key}", context.mkAuthHandler(context.unsetProjectConfigKey)).Methods("DELETE")

	r.HandleFunc("/project/{id}/key", context.mkAuthHandler(context.setKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/key", context.mkAuthHandler(context.getKey)).Methods("GET")

	r.HandleFunc("/project/{id}/crypto", context.mkAuthHandler(context.encryptData)).Methods("PUT")

	r.HandleFunc("/job", context.mkAuthHandler(context.getAllJobs)).Methods("GET")
	r.HandleFunc("/job/{id}", context.mkAuthHandler(context.getJob)).Methods("GET")
	r.HandleFunc("/job/{id}/log", context.mkAuthHandler(context.getJobLog)).Methods("GET")
	r.HandleFunc("/job/{id}/variant", context.mkAuthHandler(context.getVariants)).Methods("GET")

	r.HandleFunc("/variant/{id}", context.mkAuthHandler(context.getVariant)).Methods("GET")
	r.HandleFunc("/variant/{id}/log", context.mkAuthHandler(context.getVariantLog)).Methods("GET")
	r.HandleFunc("/variant/{id}/artifacts/{path:.*}", context.mkAuthHandler(context.getVariantArtifact)).Methods("GET")

	r.HandleFunc("/image", context.mkAuthHandler(context.getImages)).Methods("GET")
	r.HandleFunc("/image/{name:.*}", context.mkAuthHandler(context.getImage)).Methods("GET")
	r.HandleFunc("/image/{name:.*}", context.mkAuthHandler(context.setImage)).Methods("PUT")

	r.HandleFunc("/user", context.mkAuthHandler(context.getUsers)).Methods("GET")
	r.HandleFunc("/user", context.mkAuthHandler(context.createUser)).Methods("POST")
	r.HandleFunc("/user/{id}", context.mkAuthHandler(context.getUser)).Methods("GET")

	r.HandleFunc("/project/{id}/bitbucket", context.mkAuthHandler(context.startBitbucketJob)).Methods("POST")
	r.HandleFunc("/project/{id}/github", context.mkGithubAuthHandler(context.startGithubJob)).Methods("POST")

	{
		i := r.PathPrefix("/_").Subrouter()

		i.HandleFunc("/project/{id}/crypto-key", context.mkInternalApiHandler(context.getCryptoKey)).Methods("GET")
		i.HandleFunc("/job/{id}/start", context.mkInternalApiHandler(context.jobStarted)).Methods("POST")
		i.HandleFunc("/job/{id}/finish", context.mkInternalApiHandler(context.jobFinished)).Methods("POST")
		i.HandleFunc("/job/{id}/reset", context.mkInternalApiHandler(context.jobReset)).Methods("POST")
		i.HandleFunc("/job/{id}/scm", context.mkInternalApiHandler(context.addJobScmData)).Methods("PUT")
		i.HandleFunc("/variant/{id}/finish", context.mkInternalApiHandler(context.variantFinished)).Methods("POST")
		i.HandleFunc("/variant", context.mkInternalApiHandler(context.addVariant)).Methods("POST")
	}

	http.Handle("/", r)

	go func() {
		log.Infof("Starting API server on port 3000")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}()

	go func() {
		log.Infof("Starting Syslog server on port 3001")
		context.startLogServer(":3001")
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Infof("Got SIGTERM, Exiting")
}

var (
	defaultImages = map[string]string{
		"orchestration":  "bazooka/orchestration",
		"parser":         "bazooka/parser",
		"scm/fetch/git":  "bazooka/scm-git",
		"scm/fetch/hg":   "bazooka/scm-hg",
		"parser/golang":  "bazooka/parser-golang",
		"parser/go":      "bazooka/parser-golang",
		"parser/java":    "bazooka/parser-java",
		"parser/python":  "bazooka/parser-python",
		"parser/ruby":    "bazooka/parser-ruby",
		"parser/node_js": "bazooka/parser-nodejs",
		"parser/nodejs":  "bazooka/parser-nodejs",
		"parser/php":     "bazooka/parser-php",
	}
)

func ensureDefaultImagesExist(c *db.MongoConnector) error {
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

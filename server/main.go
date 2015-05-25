package main

import (
	"fmt"
	"net/http"

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
	r.HandleFunc("/project/{id}/bitbucket", context.mkAuthHandler(context.startBitbucketJob)).Methods("POST")
	r.HandleFunc("/project/{id}/job", context.mkAuthHandler(context.getJobs)).Methods("GET")

	r.HandleFunc("/project/{id}/config", context.mkAuthHandler(context.getProjectConfig)).Methods("GET")
	r.HandleFunc("/project/{id}/config/{key}", context.mkAuthHandler(context.setProjectConfigKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/config/{key}", context.mkAuthHandler(context.unsetProjectConfigKey)).Methods("DELETE")

	r.HandleFunc("/project/{id}/key", context.mkAuthHandler(context.addKey)).Methods("POST")
	r.HandleFunc("/project/{id}/key", context.mkAuthHandler(context.updateKey)).Methods("PUT")
	r.HandleFunc("/project/{id}/key", context.mkAuthHandler(context.listKeys)).Methods("GET")

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
	// r.HandleFunc("/image/{name}", context.mkAuthHandler(context.unsetImage)).Methods("DELETE")
	r.HandleFunc("/image/{name:.*}", context.mkAuthHandler(context.setImage)).Methods("PUT")

	r.HandleFunc("/user", context.mkAuthHandler(context.getUsers)).Methods("GET")
	r.HandleFunc("/user", context.mkAuthHandler(context.createUser)).Methods("POST")
	r.HandleFunc("/user/{id}", context.mkAuthHandler(context.getUser)).Methods("GET")

	r.HandleFunc("/project/{id}/github", context.mkGithubAuthHandler(context.startGithubJob)).Methods("POST")

	{
		i := r.PathPrefix("/_").Subrouter()

		i.HandleFunc("/job/{id}/status", context.mkInternalApiHandler(context.finishJob)).Methods("POST")
		i.HandleFunc("/job/{id}/scm", context.mkInternalApiHandler(context.addJobScmData)).Methods("PUT")
		i.HandleFunc("/variant/{id}/status", context.mkInternalApiHandler(context.finishVariant)).Methods("POST")
		i.HandleFunc("/variant", context.mkInternalApiHandler(context.addVariant)).Methods("POST")
	}

	http.Handle("/", r)

	go context.startLogServer()
	log.Fatal(http.ListenAndServe(":3000", nil))
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
		"ruby":           "bazooka/parser-ruby",
		"parser/node_js": "bazooka/parser-nodejs",
		"parser/nodejs":  "bazooka/parser-nodejs",
		"parser/php":     "bazooka/parser-php",
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

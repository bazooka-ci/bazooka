package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/bazooka-ci/bazooka/commons"
	"github.com/gorilla/mux"
)

type githubPayload struct {
	Ref        string       `json:"ref"`
	HeadCommit githubCommit `json:"head_commit"`
	Deleted    bool         `json:"deleted"`
}

type githubCommit struct {
	ID string `json:"id"`
}

var (
	branchRegexp = regexp.MustCompile(`refs\/heads\/(.*)`)
	tagRegexp    = regexp.MustCompile(`refs\/tags\/(.*)`)
)

func (ctx *context) mkGithubAuthHandler(f func(*request) (*response, error)) func(http.ResponseWriter, *http.Request) {
	return ctx.githubAuthenticationHandler(mkHandler(f))
}

func (ctx *context) githubAuthenticationHandler(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if ctx.githubAuth(r) {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
		}
	}
}

func (ctx *context) githubAuth(req *http.Request) bool {
	sign := req.Header.Get("X-Hub-Signature")
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body, reason is: %v\n", err)
		return false
	}
	// Replace req.Body so body can be read again
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	params := mux.Vars(req)
	project, err := ctx.Connector.GetProjectById(params["id"])
	if err != nil {
		log.Errorf("Error reading project with ID %s, reason is: %v\n", params["id"], err)
		return false
	}
	mac := hmac.New(sha1.New, []byte(project.HookKey))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal([]byte(sign), []byte(fmt.Sprintf("sha1=%x", expectedMAC)))
}

func (ctx *context) startGithubJob(r *request) (*response, error) {
	var payload githubPayload

	r.parseBody(&payload)

	if payload.Deleted {
		return noContent()
	}

	var ref string
	if branchRegexp.MatchString(payload.Ref) {
		submatch := branchRegexp.FindStringSubmatch(payload.Ref)
		if len(submatch) != 2 {
			return badRequest("Impossible to find submatch in regexp for branches")
		}
		ref = submatch[1]
	} else if tagRegexp.MatchString(payload.Ref) {
		submatch := tagRegexp.FindStringSubmatch(payload.Ref)
		if len(submatch) != 2 {
			return badRequest("Impossible to find submatch in regexp for tags")
		}
		ref = submatch[1]
	} else {
		return badRequest("ref doesn't match any know regexp for tags or branch")
	}

	return ctx.startJob(r.vars, bazooka.StartJob{
		ScmReference: ref,
	}, payload.HeadCommit.ID)
}

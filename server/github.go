package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
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

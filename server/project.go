package main

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
)

func (p *context) createProject(r *request) (*response, error) {
	var project lib.Project

	r.parseBody(&project)

	exists, err := p.Connector.HasProject(project.Name)
	switch {
	case err != nil:
		return nil, err
	case exists:
		return conflict("name is already known")
	}

	supported, err := p.Connector.HasImage(fmt.Sprintf("scm/fetch/%s", project.ScmType))
	switch {
	case err != nil:
		return nil, err
	case !supported:
		return badRequest(fmt.Sprintf("unsupported scm_type:'%s'", project.ScmType))
	}
	// TODO : validate scm_type
	// TODO : validate data by scm_type
	log.WithFields(log.Fields{
		"project": project,
	}).Info("Adding project")
	if err = p.Connector.AddProject(&project); err != nil {
		return nil, err
	}

	cryptoKey := &lib.CryptoKey{
		Content:   []byte(randSeq(32)),
		ProjectID: project.ID,
	}

	if err = p.Connector.AddCryptoKey(cryptoKey); err != nil {
		return nil, err
	}

	return created(&project, "/project/"+project.ID)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		rand.Seed(time.Now().UTC().UnixNano())
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (p *context) getProject(r *request) (*response, error) {
	project, err := p.Connector.GetProjectById(r.vars["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	return ok(&project)
}

func (p *context) getProjects(r *request) (*response, error) {
	includeStatus := len(r.r.URL.Query().Get("include-status")) > 0

	if includeStatus {
		projects, err := p.Connector.GetProjectsWithStatus()
		if err != nil {
			return nil, err
		}
		return ok(projects)
	}
	projects, err := p.Connector.GetProjects()
	if err != nil {
		return nil, err
	}
	return ok(projects)

}

func (p *context) getProjectConfig(r *request) (*response, error) {
	project, err := p.Connector.GetProjectById(r.vars["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	return ok(project.Config)
}

func (p *context) setProjectConfigKey(r *request) (*response, error) {
	id, key := r.vars["id"], r.vars["key"]
	body := r.rawBody()

	if err := p.Connector.SetProjectConfigKey(id, key, string(body)); err != nil {
		return nil, err
	}

	return noContent()
}

func (p *context) unsetProjectConfigKey(r *request) (*response, error) {
	if err := p.Connector.UnsetProjectConfigKey(r.vars["id"], r.vars["key"]); err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	return noContent()
}

func (p *context) getProjectEnv(r *request) (*response, error) {
	project, err := p.Connector.GetProjectById(r.vars["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	return ok(project.Env)
}

func (p *context) setProjectEnvKey(r *request) (*response, error) {
	id, key := r.vars["id"], r.vars["key"]
	body := r.rawBody()

	if err := p.Connector.SetProjectEnvKey(id, key, string(body)); err != nil {
		return nil, err
	}

	return noContent()
}

func (p *context) unsetProjectEnvKey(r *request) (*response, error) {
	if err := p.Connector.UnsetProjectEnvKey(r.vars["id"], r.vars["key"]); err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	return noContent()
}

package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	lib "github.com/haklop/bazooka/commons"
)

func (p *context) createProject(params map[string]string, body bodyFunc) (*response, error) {
	var project lib.Project

	body(&project)

	switch {
	case len(project.ScmURI) == 0:
		return badRequest("scm_uri is mandatory")
	case len(project.ScmType) == 0:
		return badRequest("scm_type is mandatory")
	case len(project.Name) == 0:
		return badRequest("name is mandatory")
	}

	exists, err := p.Connector.HasProject("", project.ScmType, project.ScmURI)
	switch {
	case err != nil:
		return nil, err
	case exists:
		return conflict("scm_uri is already known")
	}

	exists, err = p.Connector.HasProject(project.Name, "", "")
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

	return created(&project, "/project/"+project.ID)
}

func (p *context) getProject(params map[string]string, body bodyFunc) (*response, error) {
	project, err := p.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	return ok(&project)
}

func (p *context) getProjects(params map[string]string, body bodyFunc) (*response, error) {
	projects, err := p.Connector.GetProjects()
	log.WithFields(log.Fields{
		"projects": projects,
	}).Info("Retrieving projects")
	if err != nil {
		return nil, err
	}

	return ok(&projects)
}

package main

import lib "github.com/haklop/bazooka/commons"

func (p *context) createProject(params map[string]string, body bodyFunc) (*response, error) {
	var project lib.Project

	body(&project)

	if len(project.ScmURI) == 0 {
		return badRequest("scm_uri is mandatory")
	}

	if len(project.ScmType) == 0 {
		return badRequest("scm_type is mandatory")
	}

	existantProject, err := p.Connector.GetProject(project.ScmType, project.ScmURI)
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
	}

	if len(existantProject.ScmURI) > 0 {
		return nil, errorResponse{409, "scm_uri is already known"}
	}

	// TODO : validate scm_type
	// TODO : validate data by scm_type

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
	if err != nil {
		return nil, err
	}

	return ok(&projects)
}

package main

import lib "github.com/haklop/bazooka/commons"

func (p *Context) createProject(params map[string]string, body BodyFunc) (*Response, error) {
	var project lib.Project

	body(&project)

	if len(project.ScmURI) == 0 {
		return BadRequest("scm_uri is mandatory")
	}

	if len(project.ScmType) == 0 {
		return BadRequest("scm_type is mandatory")
	}

	existantProject, err := p.Connector.GetProject(project.ScmType, project.ScmURI)
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
	}

	if len(existantProject.ScmURI) > 0 {
		return nil, ErrorResponse{409, "scm_uri is already known"}
	}

	// TODO : validate scm_type
	// TODO : validate data by scm_type

	if err = p.Connector.AddProject(&project); err != nil {
		return nil, err
	}
	return Created(&project, "/project/"+project.ID)
}

func (p *Context) getProject(params map[string]string, body BodyFunc) (*Response, error) {
	project, err := p.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return NotFound("project not found")
	}

	return Ok(&project)
}

func (p *Context) getProjects(params map[string]string, body BodyFunc) (*Response, error) {
	projects, err := p.Connector.GetProjects()
	if err != nil {
		return nil, err
	}

	return Ok(&projects)
}

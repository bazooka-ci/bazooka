package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	lib "github.com/haklop/bazooka/commons"
)

func (p *Context) createProject(res http.ResponseWriter, req *http.Request) {
	var project lib.Project

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := decoder.Decode(&project)

	if err != nil {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "Unable to decode your json : " + err.Error(),
		})
		return
	}

	if len(project.ScmURI) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "scm_uri is mandatory",
		})

		return
	}

	if len(project.ScmType) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "scm_type is mandatory",
		})

		return
	}

	existantProject, err := p.Connector.GetProject(project.ScmType, project.ScmURI)
	if err != nil {
		if err.Error() != "not found" {
			WriteError(err, res, encoder)
			return
		}
	}

	if len(existantProject.ScmURI) > 0 {
		res.WriteHeader(409)
		encoder.Encode(&ErrorResponse{
			Code:    409,
			Message: "scm_uri is already known",
		})

		return
	}

	// TODO : validate scm_type
	// TODO : validate data by scm_type

	err = p.Connector.AddProject(&project)
	res.Header().Set("Location", "/project/"+project.ID)

	res.WriteHeader(201)
	encoder.Encode(&project)

}

func (p *Context) getProject(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	project, err := p.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "project not found",
		})

		return
	}

	res.WriteHeader(200)
	encoder.Encode(&project)
}

func (p *Context) getProjects(res http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	projects, err := p.Connector.GetProjects()
	if err != nil {
		WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&projects)
}

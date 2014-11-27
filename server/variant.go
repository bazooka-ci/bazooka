package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/haklop/bazooka/commons/mongo"
)

func (p *Context) getVariant(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	variant, err := p.Connector.GetVariantByID(params["variant_id"])
	if err != nil {
		if err.Error() != "not found" {
			WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "variant not found",
		})
		return
	}

	if params["job_id"] != variant.JobID {
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "project not found",
		})
		return
	}

	// TODO Validate project_id is correct

	res.WriteHeader(200)
	encoder.Encode(&variant)
}

func (p *Context) getVariants(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	variants, err := p.Connector.GetVariants(params["job_id"])
	if err != nil {
		WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&variants)
}

func (p *Context) getVariantLog(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	log, err := p.Connector.GetLog(&mongo.LogExample{
		ProjectID: params["id"],
		JobID:     params["job_id"],
		VariantID: params["variant_id"],
	})
	if err != nil {
		if err.Error() != "not found" {
			WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "log not found",
		})
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&log)
}

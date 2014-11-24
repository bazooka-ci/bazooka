package project

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/haklop/bazooka/server/context"
)

func (p *Handlers) getVariant(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	variant, err := p.mongoConnector.GetVariantByID(params["variant_id"])
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&context.ErrorResponse{
			Code:    404,
			Message: "variant not found",
		})
		return
	}

	if params["job_id"] != variant.JobID {
		res.WriteHeader(404)
		encoder.Encode(&context.ErrorResponse{
			Code:    404,
			Message: "project not found",
		})
		return
	}

	// TODO Validate project_id is correct

	res.WriteHeader(200)
	encoder.Encode(&variant)
}

func (p *Handlers) getVariants(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	variants, err := p.mongoConnector.GetVariants(params["job_id"])
	if err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&variants)
}

package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/haklop/bazooka/commons/mongo"
)

func (c *Context) getVariant(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	variant, err := c.Connector.GetVariantByID(params["id"])
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

	res.WriteHeader(200)
	encoder.Encode(&variant)
}

func (c *Context) getVariants(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	variants, err := c.Connector.GetVariants(params["id"])
	if err != nil {
		WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&variants)
}

func (c *Context) getVariantLog(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	log, err := c.Connector.GetLog(&mongo.LogExample{
		VariantID: params["id"],
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

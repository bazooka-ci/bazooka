package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	lib "github.com/haklop/bazooka/commons"
)

func (p *Context) createFetcher(res http.ResponseWriter, req *http.Request) {
	var fetcher lib.ScmFetcher

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := decoder.Decode(&fetcher)

	if err != nil {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "Unable to decode your json : " + err.Error(),
		})
		return
	}

	if len(fetcher.Name) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "name is mandatory",
		})

		return
	}

	if len(fetcher.ImageName) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "image_name is mandatory",
		})

		return
	}

	existantFetcher, err := p.Connector.GetFetcherByName(fetcher.Name)
	if err != nil {
		if err.Error() != "not found" {
			WriteError(err, res, encoder)
			return
		}
	}

	if len(existantFetcher.Name) > 0 {
		res.WriteHeader(409)
		encoder.Encode(&ErrorResponse{
			Code:    409,
			Message: "name is already known",
		})

		return
	}

	err = p.Connector.AddFetcher(&fetcher)
	res.Header().Set("Location", "/fetcher/"+fetcher.ID)

	res.WriteHeader(201)
	encoder.Encode(&fetcher)

}

func (p *Context) getFetcher(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	fetcher, err := p.Connector.GetFetcherById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			WriteError(err, res, encoder)
			return
		} else {
			res.WriteHeader(404)
			encoder.Encode(&ErrorResponse{
				Code:    404,
				Message: "fetcher not found",
			})

			return
		}
	}

	res.WriteHeader(200)
	encoder.Encode(&fetcher)
}

func (p *Context) getFetchers(res http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	fetchers, err := p.Connector.GetFetchers()
	if err != nil {
		WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&fetchers)
}

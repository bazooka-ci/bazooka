package fetcher

import (
	"encoding/json"
	"net/http"

	lib "github.com/haklop/bazooka/commons"
	"github.com/haklop/bazooka/commons/mongo"
	"github.com/gorilla/mux"
	"github.com/haklop/bazooka/server/context"
)

type Handlers struct {
	mongoConnector *mongo.MongoConnector
	env            map[string]string
}

func (p *Handlers) SetHandlers(r *mux.Router, serverContext context.Context) {
	p.mongoConnector = serverContext.Connector
	p.env = serverContext.Env

	r.HandleFunc("/", p.createFetcher).Methods("POST")
	r.HandleFunc("/", p.getFetchers).Methods("GET")
	r.HandleFunc("/{id}", p.getFetcher).Methods("GET")
}

func (p *Handlers) createFetcher(res http.ResponseWriter, req *http.Request) {
	var fetcher lib.ScmFetcher

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := decoder.Decode(&fetcher)

	if err != nil {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "Unable to decode your json : " + err.Error(),
		})
		return
	}

	if len(fetcher.Name) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "name is mandatory",
		})

		return
	}

	if len(fetcher.ImageName) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "image_name is mandatory",
		})

		return
	}

	existantFetcher, err := p.mongoConnector.GetFetcherByName(fetcher.Name)
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
	}

	if len(existantFetcher.Name) > 0 {
		res.WriteHeader(409)
		encoder.Encode(&context.ErrorResponse{
			Code:    409,
			Message: "name is already known",
		})

		return
	}

	err = p.mongoConnector.AddFetcher(&fetcher)
	res.Header().Set("Location", "/fetcher/"+fetcher.ID)

	res.WriteHeader(201)
	encoder.Encode(&fetcher)

}

func (p *Handlers) getFetcher(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	fetcher, err := p.mongoConnector.GetFetcherById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		} else {
			res.WriteHeader(404)
			encoder.Encode(&context.ErrorResponse{
				Code:    404,
				Message: "fetcher not found",
			})

			return
		}
	}

	res.WriteHeader(200)
	encoder.Encode(&fetcher)
}

func (p *Handlers) getFetchers(res http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	fetchers, err := p.mongoConnector.GetFetchers()
	if err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&fetchers)
}

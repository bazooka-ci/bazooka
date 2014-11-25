package project

import (
	"encoding/json"
	"net/http"

	lib "github.com/haklop/bazooka/commons"

	"github.com/gorilla/mux"
	"github.com/haklop/bazooka/server/context"
)

func (p *Handlers) getJobLog(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	log, err := p.mongoConnector.GetLog(&lib.LogEntry{
		JobID: params["job_id"],
	})
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&context.ErrorResponse{
			Code:    404,
			Message: "log not found",
		})
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&log)
}

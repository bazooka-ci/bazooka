package context

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerHandler interface {
	SetHandlers(router *mux.Router, context Context)
}

func WriteError(err error, res http.ResponseWriter, encoder *json.Encoder) {
	res.WriteHeader(500)
	encoder.Encode(&ErrorResponse{
		Code:    500,
		Message: err.Error(),
	})
}

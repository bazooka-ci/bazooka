package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/haklop/bazooka/commons/mongo"
)

const (
	BazookaEnvSCMKeyfile = "BZK_SCM_KEYFILE"
	BazookaEnvHome       = "BZK_HOME"
	BazookaEnvDockerSock = "BZK_DOCKERSOCK"
	BazookaEnvMongoAddr  = "MONGO_PORT_27017_TCP_ADDR"
	BazookaEnvMongoPort  = "MONGO_PORT_27017_TCP_PORT"

	DockerSock     = "/var/run/docker.sock"
	DockerEndpoint = "unix://" + DockerSock
	BazookaHome    = "/bazooka"
)

type ErrorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

type Context struct {
	Connector      *mongo.MongoConnector
	DockerEndpoint string
	Env            map[string]string
}

func WriteError(err error, res http.ResponseWriter, encoder *json.Encoder) {
	res.WriteHeader(500)
	encoder.Encode(&ErrorResponse{
		Code:    500,
		Message: err.Error(),
	})
}

type BodyFunc func(interface{})

type Response struct {
	Code    int
	Payload interface{}
	Headers map[string]string
}

func Ok(payload interface{}) (*Response, error) {
	return &Response{
		Code:    200,
		Payload: payload,
	}, nil
}

func Created(payload interface{}, location string) (*Response, error) {
	return &Response{
		201,
		payload,
		map[string]string{"Location": location},
	}, nil
}
func Accepted(payload interface{}, location string) (*Response, error) {
	return &Response{
		202,
		payload,
		map[string]string{"Location": location},
	}, nil
}
func BadRequest(msg string) (*Response, error) {
	return nil, &ErrorResponse{400, msg}
}

func NotFound(msg string) (*Response, error) {
	return nil, &ErrorResponse{404, msg}
}

func MkHandler(f func(map[string]string, BodyFunc) (*Response, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bf := func(b interface{}) {
			defer r.Body.Close()
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(b); err != nil {
				panic(ErrorResponse{400, "Unable to decode your json : " + err.Error()})
			}
		}

		encoder := json.NewEncoder(w)

		dispatchError := func(err error) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			switch e := err.(type) {
			case ErrorResponse:
				w.WriteHeader(400)
				encoder.Encode(e)
			default:
				WriteError(e, w, encoder)
			}
		}

		defer func() {
			if r := recover(); r != nil {
				switch rt := r.(type) {
				case error:
					dispatchError(rt)
				default:
					WriteError(fmt.Errorf("Caught a panic: %v", r), w, encoder)
				}
			}
		}()

		rb, err := f(mux.Vars(r), bf)

		if err != nil {
			dispatchError(err)
			return
		}

		if rb != nil {
			fmt.Printf("rb=%#v\n", rb)

			for k, v := range rb.Headers {
				fmt.Printf("add header %s=%s\n", k, v)
				w.Header().Set(k, v)
			}

			if rb.Payload != nil {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(rb.Code)
				encoder.Encode(&rb.Payload)
			} else {
				w.WriteHeader(rb.Code)
			}
		}

	}
}

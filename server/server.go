package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bazooka-ci/bazooka/commons/mongo"
	basic "github.com/haklop/go-http-basic-auth"
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

type errorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e errorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

type context struct {
	Connector      *mongo.MongoConnector
	DockerEndpoint string
	Env            map[string]string
}

func writeError(err error, res http.ResponseWriter) {
	res.WriteHeader(500)
	json.NewEncoder(res).Encode(&errorResponse{
		Code:    500,
		Message: err.Error(),
	})
}

type bodyFunc func(interface{})

type response struct {
	Code    int
	Payload interface{}
	Headers map[string]string
}

func ok(payload interface{}) (*response, error) {
	return &response{
		Code:    200,
		Payload: payload,
	}, nil
}

func noContent() (*response, error) {
	return &response{
		Code: 204,
	}, nil
}

func created(payload interface{}, location string) (*response, error) {
	return &response{
		201,
		payload,
		map[string]string{"Location": location},
	}, nil
}
func accepted(payload interface{}, location string) (*response, error) {
	return &response{
		202,
		payload,
		map[string]string{"Location": location},
	}, nil
}
func badRequest(msg string) (*response, error) {
	return nil, &errorResponse{400, msg}
}

func notFound(msg string) (*response, error) {
	return nil, &errorResponse{404, msg}
}

func conflict(msg string) (*response, error) {
	return nil, &errorResponse{409, msg}
}

func unauthorized() (*response, error) {
	return nil, &errorResponse{401, "Unauthorized"}
}

func mkHandler(f func(map[string]string, bodyFunc) (*response, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		bf := func(b interface{}) {
			defer r.Body.Close()
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(b); err != nil {
				panic(errorResponse{400, "Unable to decode your json : " + err.Error()})
			}
		}

		encoder := json.NewEncoder(w)

		dispatchError := func(err error) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			switch e := err.(type) {
			case errorResponse:
				w.WriteHeader(e.Code)
				encoder.Encode(e)
			case *errorResponse:
				w.WriteHeader(e.Code)
				encoder.Encode(e)
			default:
				writeError(e, w)
			}
		}

		defer func() {
			if r := recover(); r != nil {
				switch rt := r.(type) {
				case error:
					dispatchError(rt)
				default:
					writeError(fmt.Errorf("Caught a panic: %v", r), w)
				}
			}
		}()

		rb, err := f(mux.Vars(r), bf)

		if err != nil {
			dispatchError(err)
			return
		}

		if rb != nil {
			for k, v := range rb.Headers {
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

func (ctx *context) authenticationHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := ctx.Connector.GetUsers()
		if err != nil {
			// TODO error
			log.Fatal(err)
		}

		if len(users) > 0 {
			authenticator := basic.NewAuthenticator(ctx.userAuthentication, "bazooka")

			authenticator.Wrap(next).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})

}

func (ctx *context) userAuthentication(email string, password string) bool {
	return ctx.Connector.ComparePassword(email, password)
}

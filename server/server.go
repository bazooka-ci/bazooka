package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	basic "github.com/haklop/go-http-basic-auth"
	validator "gopkg.in/bluesuncorp/validator.v5"
)

type errorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e errorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func writeError(err error, res http.ResponseWriter) {
	res.WriteHeader(500)
	json.NewEncoder(res).Encode(&errorResponse{
		Code:    500,
		Message: err.Error(),
	})
}

func flushResponse(w http.ResponseWriter) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

type request struct {
	w        http.ResponseWriter
	r        *http.Request
	vars     map[string]string
	validate *validator.Validate
}

func (r *request) parseBody(into interface{}) {
	defer r.r.Body.Close()
	decoder := json.NewDecoder(r.r.Body)
	if err := decoder.Decode(into); err != nil {
		panic(errorResponse{400, "Unable to decode your json : " + err.Error()})
	}
	if err := r.validate.Struct(into); err != nil {
		for k, v := range err.Errors {
			switch v.Tag {
			case "required":
				pointerValue := reflect.ValueOf(into)
				structType := reflect.TypeOf(pointerValue.Elem().Interface())

				if fieldData, ok := structType.FieldByName(k); ok {
					panic(errorResponse{400, fmt.Sprintf("%s is required", fieldData.Tag.Get("json"))})
				} else {
					panic(errorResponse{400, fmt.Sprintf("%s is required", v.Field)})
				}

			default:
				panic(errorResponse{400, v.Error()})
			}

		}
	}

}

func (r *request) rawBody() []byte {
	defer r.r.Body.Close()
	body, err := ioutil.ReadAll(r.r.Body)
	if err != nil {
		panic(errorResponse{400, "Unable to read request payload : " + err.Error()})
	}
	return body
}

func (r *request) query(key string) string {
	return r.r.URL.Query().Get(key)
}

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

func (ctx *context) mkAuthHandler(f func(*request) (*response, error)) func(http.ResponseWriter, *http.Request) {
	return ctx.authenticationHandler(mkHandler(f))
}

func mkHandler(f func(*request) (*response, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New("validate", validator.BakedInValidators)

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

		wrapped := &request{
			w:        w,
			r:        r,
			vars:     mux.Vars(r),
			validate: validate,
		}

		rb, err := f(wrapped)

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
	})
}

func (ctx *context) authenticationHandler(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := ctx.connector.GetUsers()
		if err != nil {
			// TODO handle error properly
			log.Fatal(err)
		}

		if len(users) > 0 {
			authenticator := basic.NewAuthenticator(ctx.userAuthentication, "bazooka")
			authenticator.Wrap(next).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}
}

func (ctx *context) mkInternalApiHandler(f func(*request) (*response, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mkHandler(f).ServeHTTP(w, r)
	}
}

func (ctx *context) userAuthentication(email string, password string) bool {
	return ctx.connector.ComparePassword(email, password)
}

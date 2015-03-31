package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bazooka-ci/bazooka/commons/mongo"
)

func (c *context) getVariant(params map[string]string, body bodyFunc) (*response, error) {
	variant, err := c.Connector.GetVariantByID(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("variant not found")
	}

	return ok(&variant)
}

func (c *context) getVariants(params map[string]string, body bodyFunc) (*response, error) {
	variants, err := c.Connector.GetVariants(params["id"])
	if err != nil {
		return nil, err
	}

	return ok(&variants)
}

func (c *context) getVariantLog(params map[string]string, body bodyFunc) (*response, error) {
	log, err := c.Connector.GetLog(&mongo.LogExample{
		VariantID: params["id"],
	})
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("log not found")
	}

	return ok(&log)
}

func (c *context) getVariantArtifacts(w http.ResponseWriter, r *http.Request) {
	vid := mux.Vars(r)["id"]
	variant, err := c.Connector.GetVariantByID(vid)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)

		if err.Error() != "not found" {
			w.WriteHeader(500)
			encoder.Encode(err)
			return
		}

		w.WriteHeader(404)
		encoder.Encode(fmt.Errorf("variant not found"))
		return
	}

	buildFolder := fmt.Sprintf("/bazooka/build/%s/%s/artifacts/%s", variant.ProjectID, variant.JobID, variant.ID)
	prefix := fmt.Sprintf("/variant/%s/artifacts/", vid)
	http.StripPrefix(prefix, http.FileServer(http.Dir(buildFolder))).ServeHTTP(w, r)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/bazooka-ci/bazooka/commons/mongo"

	log "github.com/Sirupsen/logrus"

	lib "github.com/bazooka-ci/bazooka/commons"
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

func (c *context) getVariantLog(w http.ResponseWriter, r *http.Request) {
	follow := len(r.URL.Query().Get("follow")) > 0

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

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	logOutput := json.NewEncoder(w)

	query := &mongo.LogExample{
		VariantID: vid,
	}

	logs, err := c.Connector.GetLog(query)
	if !follow {
		logOutput.Encode(logs)
		return
	}

	for _, l := range logs {
		logOutput.Encode(l)
	}
	flushResponse(w)
	lastTime := variantLastLogTime(variant, logs)

	for {
		time.Sleep(1000 * time.Millisecond)
		query.After = lastTime
		logs, err := c.Connector.GetLog(query)
		if err != nil {
			log.Errorf("Error while retrieving logs: %v", err)
			return
		}
		if len(logs) > 0 {
			lastTime = variantLastLogTime(variant, logs)
			for _, l := range logs {
				logOutput.Encode(l)
			}
			flushResponse(w)
		}
		variant, err := c.Connector.GetVariantByID(vid)
		if err != nil {
			log.Errorf("Error while retrieving variant: %v", err)
			return
		}
		if variant.Status != lib.JOB_RUNNING {
			return
		}
	}
}

func variantLastLogTime(variant *lib.Variant, logs []lib.LogEntry) time.Time {
	if len(logs) == 0 {
		return variant.Started
	}
	return logs[len(logs)-1].Time
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

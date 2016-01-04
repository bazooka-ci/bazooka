package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bazooka-ci/bazooka/server/db"

	log "github.com/Sirupsen/logrus"

	lib "github.com/bazooka-ci/bazooka/commons"
)

func (c *context) getVariant(r *request) (*response, error) {
	variant, err := c.connector.GetVariantByID(r.vars["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("variant not found")
	}

	return ok(&variant)
}

func (c *context) getVariants(r *request) (*response, error) {
	variants, err := c.connector.GetVariants(r.vars["id"])
	if err != nil {
		return nil, err
	}

	return ok(&variants)
}

func (c *context) getVariantLog(r *request) (*response, error) {
	follow := len(r.query("follow")) > 0
	strictJson := len(r.query("strict-json")) > 0

	vid := r.vars["id"]

	variant, err := c.connector.GetVariantByID(vid)

	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("variant not found")
	}

	w := r.w
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	logOutput := json.NewEncoder(w)

	query := &db.LogExample{
		VariantID: vid,
	}

	logs, err := c.connector.GetLog(query)
	if !follow {
		logOutput.Encode(logs)
		return nil, nil
	}

	if strictJson {
		w.Write([]byte("["))
		defer w.Write([]byte("]"))
	}

	writtenEntries := 0
	for _, l := range logs {
		if writtenEntries > 0 && strictJson {
			w.Write([]byte(","))
		}
		logOutput.Encode(l)
		writtenEntries++
	}
	flushResponse(w)

	if variant.Status != lib.JOB_RUNNING {
		return nil, nil
	}

	lastTime := variantLastLogTime(variant, logs)

	for {
		time.Sleep(1000 * time.Millisecond)
		query.After = lastTime
		logs, err := c.connector.GetLog(query)
		if err != nil {
			log.Errorf("Error while retrieving logs: %v", err)
			return nil, nil
		}
		if len(logs) > 0 {
			lastTime = variantLastLogTime(variant, logs)
			for _, l := range logs {
				if writtenEntries > 0 && strictJson {
					w.Write([]byte(","))
				}
				logOutput.Encode(l)
				writtenEntries++
			}
			flushResponse(w)
		}
		variant, err := c.connector.GetVariantByID(vid)
		if err != nil {
			log.Errorf("Error while retrieving variant: %v", err)
			return nil, nil
		}
		if variant.Status != lib.JOB_RUNNING {
			return nil, nil
		}
	}
}

func variantLastLogTime(variant *lib.Variant, logs []lib.LogEntry) time.Time {
	if len(logs) == 0 {
		return variant.Started
	}
	return logs[len(logs)-1].Time
}

func (c *context) getVariantArtifact(r *request) (*response, error) {
	vid := r.vars["id"]

	variant, err := c.connector.GetVariantByID(vid)

	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("variant not found")
	}

	buildFolder := fmt.Sprintf("/bazooka/build/%s/%s/artifacts/%s", variant.ProjectID, variant.JobID, variant.ID)
	prefix := fmt.Sprintf("/variant/%s/artifacts/", vid)
	http.StripPrefix(prefix, http.FileServer(http.Dir(buildFolder))).ServeHTTP(r.w, r.r)
	return nil, nil
}

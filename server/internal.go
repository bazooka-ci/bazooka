package main

import (
	"fmt"

	lib "github.com/bazooka-ci/bazooka/commons"

	"time"
)

func (c *context) finishJob(r *request) (*response, error) {
	var f lib.FinishData
	r.parseBody(&f)
	if f.Time.IsZero() {
		f.Time = time.Now()
	}
	if err := c.connector.FinishJob(r.vars["id"], f.Status, f.Time); err != nil {
		return nil, err
	}

	return noContent()
}

func (c *context) finishVariant(r *request) (*response, error) {
	var f lib.FinishData
	r.parseBody(&f)
	if f.Time.IsZero() {
		f.Time = time.Now()
	}
	if err := c.connector.FinishVariant(r.vars["id"], f.Status, f.Time, f.Artifacts); err != nil {
		return nil, err
	}

	return noContent()
}

func (c *context) addJobScmData(r *request) (*response, error) {
	var m lib.SCMMetadata
	r.parseBody(&m)
	if err := c.connector.AddJobSCMMetadata(r.vars["id"], &m); err != nil {
		return nil, err
	}

	return noContent()
}

func (c *context) addVariant(r *request) (*response, error) {
	var variant lib.Variant
	r.parseBody(&variant)

	if err := c.connector.AddVariant(&variant); err != nil {
		return nil, err
	}

	return created(&variant, fmt.Sprintf("/variant/%s", variant.ID))
}

func (c *context) getCryptoKey(r *request) (*response, error) {
	key, err := c.connector.GetProjectCryptoKey(r.vars["id"])
	if err != nil {
		return nil, err
	}

	return ok(&key)
}

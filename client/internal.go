package client

import (
	"fmt"
	"net/url"
	"time"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/racker/perigee"
)

type Internal struct {
	config *Config
}

func (in *Internal) FinishJob(jobID string, status lib.JobStatus) error {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/job/%s/status", url.QueryEscape(jobID)))
	if err != nil {
		return err
	}

	return perigee.Post(requestURL, perigee.Options{
		ReqBody: lib.FinishData{
			Status: status,
		},
		OkCodes:    []int{204},
		SetHeaders: in.config.authenticateRequest,
	})
}

func (in *Internal) FinishVariant(variantID string, status lib.JobStatus, when time.Time, artifacts []string) error {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/variant/%s/status", url.QueryEscape(variantID)))
	if err != nil {
		return err
	}

	return perigee.Post(requestURL, perigee.Options{
		ReqBody: lib.FinishData{
			Status:    status,
			Time:      when,
			Artifacts: artifacts,
		},
		OkCodes:    []int{204},
		SetHeaders: in.config.authenticateRequest,
	})
}

func (in *Internal) AddJobSCMMetadata(jobID string, m *lib.SCMMetadata) error {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/job/%s/scm", url.QueryEscape(jobID)))
	if err != nil {
		return err
	}

	return perigee.Put(requestURL, perigee.Options{
		ReqBody:    &m,
		OkCodes:    []int{204},
		SetHeaders: in.config.authenticateRequest,
	})
}

func (in *Internal) AddVariant(variant *lib.Variant) (*lib.Variant, error) {
	requestURL, err := in.config.getRequestURL("_/variant")
	if err != nil {
		return nil, err
	}
	var createdVariant lib.Variant
	err = perigee.Post(requestURL, perigee.Options{
		ReqBody:    &variant,
		Results:    &createdVariant,
		OkCodes:    []int{201},
		SetHeaders: in.config.authenticateRequest,
	})
	return &createdVariant, err
}

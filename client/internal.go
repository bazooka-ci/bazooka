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

func (in *Internal) GetProjectCryptoKey(projectID string) (*lib.CryptoKey, error) {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/project/%s/crypto-key", projectID))
	if err != nil {
		return nil, err
	}
	var cryptoKey lib.CryptoKey
	err = perigee.Get(requestURL, perigee.Options{
		Results:    &cryptoKey,
		OkCodes:    []int{200},
		SetHeaders: in.config.authenticateRequest,
	})
	return &cryptoKey, err
}

func (in *Internal) MarkJobAsStarted(jobID string) error {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/job/%s/start", url.QueryEscape(jobID)))
	if err != nil {
		return err
	}

	return perigee.Post(requestURL, perigee.Options{
		ReqBody: lib.FinishData{
			Status: lib.JOB_RUNNING,
		},
		OkCodes:    []int{204},
		SetHeaders: in.config.authenticateRequest,
	})
}

func (in *Internal) MarkJobAsFinished(jobID string, status lib.JobStatus) error {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/job/%s/finish", url.QueryEscape(jobID)))
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

func (in *Internal) ResetJob(jobID string) error {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/job/%s/reset", url.QueryEscape(jobID)))
	if err != nil {
		return err
	}

	return perigee.Post(requestURL, perigee.Options{
		OkCodes:    []int{204},
		SetHeaders: in.config.authenticateRequest,
	})
}

func (in *Internal) MarkVariantAsFinished(variantID string, status lib.JobStatus, when time.Time, artifacts []string) error {
	requestURL, err := in.config.getRequestURL(fmt.Sprintf("_/variant/%s/finish", url.QueryEscape(variantID)))
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

func (in *Internal) Heartbeat() {

}

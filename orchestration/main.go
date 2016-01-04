package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

func main() {
	// TODO add validation
	start := time.Now()

	context := initContext()

	log.WithFields(log.Fields{
		"environment": context,
	}).Info("Starting Orchestration")

	f := &SCMFetcher{
		context: context,
	}
	err := f.Fetch()
	if err != nil && context.reuseScm {
		log.Info("First SCM fetch with bzk.scm.reuse true failed, retrying with a clean SCM fetch")
		f.update = false
		err = f.Fetch()
	}
	if err != nil {
		clientErr := context.client.Internal.MarkJobAsFinished(context.jobID, lib.JOB_ERRORED)
		if clientErr != nil {
			log.Fatal(err, clientErr)
		}
		log.Fatal(err)
	}

	p := &Parser{
		context: context,
	}
	parsedVariants, err := p.Parse()
	if err != nil {
		clientErr := context.client.Internal.MarkJobAsFinished(context.jobID, lib.JOB_ERRORED)
		if clientErr != nil {
			log.Fatal(err, clientErr)
		}
		log.Fatal(err)
	}

	for i, v := range parsedVariants {
		variant := &lib.Variant{
			Started:   time.Now(),
			Status:    lib.JOB_RUNNING,
			Number:    i,
			ProjectID: context.projectID,
			JobID:     context.jobID,
			Metas:     v.meta,
		}
		var err error
		variant, err = context.client.Internal.AddVariant(variant)
		if err != nil {
			clientErr := context.client.Internal.MarkJobAsFinished(context.jobID, lib.JOB_ERRORED)
			if clientErr != nil {
				log.Fatal(err, clientErr)
			}
			log.Fatal(err)
		}
		v.variant = variant
	}

	b := &Builder{
		context:  context,
		variants: parsedVariants,
	}

	if err := b.Build(); err != nil {
		clientErr := context.client.Internal.MarkJobAsFinished(context.jobID, lib.JOB_ERRORED)
		if clientErr != nil {
			log.Fatal(err, clientErr)
		}
		log.Fatal(err)
	}

	// variantsToBuild are the variants that we succeeded in generating a doocker image for them
	variantsToBuild := []*variantData{}
	for _, vd := range parsedVariants {
		switch vd.variant.Status {
		case lib.JOB_ERRORED:
			if err := context.client.Internal.MarkVariantAsFinished(vd.variant.ID, lib.JOB_ERRORED, vd.variant.Completed, nil); err != nil {
				log.Fatal(err)
			}
		default:
			variantsToBuild = append(variantsToBuild, vd)
		}
	}

	r := &Runner{
		variants: variantsToBuild,
		context:  context,
	}

	err = r.Run()
	if err != nil {
		clientErr := context.client.Internal.MarkJobAsFinished(context.jobID, lib.JOB_ERRORED)
		if clientErr != nil {
			log.Fatal(err, clientErr)
		}
		log.Fatal(err)
	}

	for _, vd := range variantsToBuild {
		if err := context.client.Internal.MarkVariantAsFinished(vd.variant.ID, vd.variant.Status, vd.variant.Completed, vd.variant.Artifacts); err != nil {
			log.Fatal(err)
		}
	}

	var (
		errorCount   = 0
		successCount = 0
		failCount    = 0
	)
	for _, vd := range parsedVariants {
		switch vd.variant.Status {
		case lib.JOB_ERRORED:
			errorCount++
		case lib.JOB_SUCCESS:
			successCount++
		case lib.JOB_FAILED:
			failCount++
		default:
			log.Fatal(fmt.Errorf("Found a variant without a status %v", vd))
		}
	}

	log.WithFields(log.Fields{
		"ERRORED":   strconv.Itoa(errorCount),
		"SUCCEEDED": strconv.Itoa(successCount),
		"FAILED":    strconv.Itoa(failCount),
	}).Info("Job Completed")

	var jobStatus lib.JobStatus
	switch {
	case errorCount > 0:
		jobStatus = lib.JOB_ERRORED
	case failCount > 0:
		jobStatus = lib.JOB_FAILED
	default:
		jobStatus = lib.JOB_SUCCESS
	}

	if err = context.client.Internal.MarkJobAsFinished(context.jobID, jobStatus); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	log.WithFields(log.Fields{
		"elapsed": elapsed,
	}).Info("Job Orchestration finished")
}

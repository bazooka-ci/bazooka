package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"
	docker "github.com/bywan/go-dockercommand"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

type Logger func(image string, variant string, container *docker.Container)

func main() {
	// TODO add validation
	start := time.Now()

	context := initContext()
	defer context.cleanup()

	var containerLogger Logger = func(image string, variantID string, container *docker.Container) {
		r, w := io.Pipe()
		container.StreamLogs(w)
		context.connector.FeedLog(r, lib.LogEntry{
			ProjectID: context.projectID,
			JobID:     context.jobID,
			VariantID: variantID,
			Image:     image,
		})
	}

	//redirect the log to mongo
	func() {
		r, w := io.Pipe()
		log.SetOutput(io.MultiWriter(os.Stdout, w))
		context.connector.FeedLog(r, lib.LogEntry{
			ProjectID: context.projectID,
			JobID:     context.jobID,
			Image:     "bazooka/orchestration",
		})
	}()

	log.WithFields(log.Fields{
		"environment": context,
	}).Info("Starting Orchestration")

	f := &SCMFetcher{
		context: context,
	}
	err := f.Fetch(containerLogger)
	if err != nil && context.reuseScm {
		log.Info("First SCM fetch with bzk.scm.reuse true failed, retrying with a clean SCM fetch")
		f.update = false
		err = f.Fetch(containerLogger)
	}
	if err != nil {
		mongoErr := context.connector.FinishJob(context.jobID, lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			log.Fatal(mongoErr)
		}
		log.Fatal(err)
	}

	p := &Parser{
		context: context,
	}
	parsedVariants, err := p.Parse(containerLogger)
	if err != nil {
		mongoErr := context.connector.FinishJob(context.jobID, lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			log.Fatal(err, mongoErr)
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
		err := context.connector.AddVariant(variant)
		if err != nil {
			mongoErr := context.connector.FinishJob(context.jobID, lib.JOB_ERRORED, time.Now())
			if mongoErr != nil {
				log.Fatal(err, mongoErr)
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
		mongoErr := context.connector.FinishJob(context.jobID, lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			log.Fatal(mongoErr)
		}
		log.Fatal(err)
	}

	// variantsToBuild are the variants that we succeeded in generating a doocker image for them
	variantsToBuild := []*variantData{}
	for _, vd := range parsedVariants {
		switch vd.variant.Status {
		case lib.JOB_ERRORED:
			if err := context.connector.FinishVariant(vd.variant.ID, lib.JOB_ERRORED, vd.variant.Completed, nil); err != nil {
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

	err = r.Run(containerLogger)
	if err != nil {
		mongoErr := context.connector.FinishJob(context.jobID, lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			log.Fatal(mongoErr)
		}
		log.Fatal(err)
	}

	for _, vd := range variantsToBuild {
		if err := context.connector.FinishVariant(vd.variant.ID, vd.variant.Status, vd.variant.Completed, vd.variant.Artifacts); err != nil {
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

	if err = context.connector.FinishJob(context.jobID, jobStatus, time.Now()); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	log.WithFields(log.Fields{
		"elapsed": elapsed,
	}).Info("Job Orchestration finished")
}

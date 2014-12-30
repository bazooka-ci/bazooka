package main

import (
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"
	lib "github.com/haklop/bazooka/commons"
	bzklog "github.com/haklop/bazooka/commons/logs"
	"github.com/haklop/bazooka/commons/mongo"
)

const (
	CheckoutFolderPattern = "%s/source"
	WorkdirFolderPattern  = "%s/work"
	MetaFolderPattern     = "%s/meta"

	BazookaInput   = "/bazooka"
	DockerSock     = "/var/run/docker.sock"
	DockerEndpoint = "unix://" + DockerSock

	BazookaEnvSCM          = "BZK_SCM"
	BazookaEnvSCMUrl       = "BZK_SCM_URL"
	BazookaEnvSCMReference = "BZK_SCM_REFERENCE"
	BazookaEnvSCMKeyfile   = "BZK_SCM_KEYFILE"
	BazookaEnvProjectID    = "BZK_PROJECT_ID"
	BazookaEnvJobID        = "BZK_JOB_ID"
	BazookaEnvHome         = "BZK_HOME"
	BazookaEnvDockerSock   = "BZK_DOCKERSOCK"
	BazookaEnvMongoAddr    = "MONGO_PORT_27017_TCP_ADDR"
	BazookaEnvMongoPort    = "MONGO_PORT_27017_TCP_PORT"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
}

type Logger func(image string, variant string, container *docker.Container)

func main() {
	// TODO add validation
	start := time.Now()

	// Configure Mongo
	connector := mongo.NewConnector()
	defer connector.Close()

	env := map[string]string{
		BazookaEnvSCM:          os.Getenv(BazookaEnvSCM),
		BazookaEnvSCMUrl:       os.Getenv(BazookaEnvSCMUrl),
		BazookaEnvSCMReference: os.Getenv(BazookaEnvSCMReference),
		BazookaEnvSCMKeyfile:   os.Getenv(BazookaEnvSCMKeyfile),
		BazookaEnvProjectID:    os.Getenv(BazookaEnvProjectID),
		BazookaEnvJobID:        os.Getenv(BazookaEnvJobID),
		BazookaEnvHome:         os.Getenv(BazookaEnvHome),
		BazookaEnvDockerSock:   os.Getenv(BazookaEnvDockerSock),
	}

	var containerLogger Logger = func(image string, variantID string, container *docker.Container) {
		r, w := io.Pipe()
		container.StreamLogs(w)
		connector.FeedLog(r, lib.LogEntry{
			ProjectID: env[BazookaEnvProjectID],
			JobID:     env[BazookaEnvJobID],
			VariantID: variantID,
			Image:     image,
		})
	}

	//redirect the log to mongo
	func() {
		r, w := io.Pipe()
		log.SetOutput(io.MultiWriter(os.Stdout, w))
		connector.FeedLog(r, lib.LogEntry{
			ProjectID: env[BazookaEnvProjectID],
			JobID:     env[BazookaEnvJobID],
			Image:     "bazooka/orchestration",
		})
	}()

	log.WithFields(log.Fields{
		"environment": env,
	}).Info("Starting Orchestration")

	checkoutFolder := fmt.Sprintf(CheckoutFolderPattern, env[BazookaEnvHome])
	metaFolder := fmt.Sprintf(MetaFolderPattern, env[BazookaEnvHome])
	f := &SCMFetcher{
		MongoConnector: connector,
		Options: &FetchOptions{
			Scm:         env[BazookaEnvSCM],
			URL:         env[BazookaEnvSCMUrl],
			Reference:   env[BazookaEnvSCMReference],
			JobID:       env[BazookaEnvJobID],
			LocalFolder: checkoutFolder,
			MetaFolder:  metaFolder,
			KeyFile:     env[BazookaEnvSCMKeyfile],
			Env:         env,
		},
	}
	if err := f.Fetch(containerLogger); err != nil {
		mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			log.Fatal(mongoErr)
		}
		log.Fatal(err)
	}

	p := &Parser{
		MongoConnector: connector,
		Options: &ParseOptions{
			InputFolder:    checkoutFolder,
			OutputFolder:   fmt.Sprintf(WorkdirFolderPattern, env[BazookaEnvHome]),
			DockerSock:     env[BazookaEnvDockerSock],
			HostBaseFolder: checkoutFolder,
			MetaFolder:     metaFolder,
			Env:            env,
		},
	}
	parsedVariants, err := p.Parse(containerLogger)
	if err != nil {
		mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			log.Fatal(err, mongoErr)
		}
		log.Fatal(err)
	}

	for i, v := range parsedVariants {
		variant := &lib.Variant{
			Started: time.Now(),
			Status:  lib.JOB_RUNNING,
			Number:  i,
			JobID:   env[BazookaEnvJobID],
			Metas:   v.meta,
		}
		err := connector.AddVariant(variant)
		if err != nil {
			mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
			if mongoErr != nil {
				log.Fatal(err, mongoErr)
			}
			log.Fatal(err)
		}
		v.variant = variant

	}

	b := &Builder{
		Options: &BuildOptions{
			SourceFolder: fmt.Sprintf(CheckoutFolderPattern, BazookaInput),
			ProjectID:    env[BazookaEnvProjectID],
			Variants:     parsedVariants,
		},
	}

	if err := b.Build(); err != nil {
		mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
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
			if err := connector.FinishVariant(vd.variant.ID, lib.JOB_ERRORED, vd.variant.Completed); err != nil {
				log.Fatal(err)
			}
		default:
			variantsToBuild = append(variantsToBuild, vd)
		}
	}

	r := &Runner{
		Variants: variantsToBuild,
		Env:      env,
		Mongo:    connector,
	}

	err = r.Run(containerLogger)
	if err != nil {
		mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			log.Fatal(mongoErr)
		}
		log.Fatal(err)
	}

	for _, vd := range variantsToBuild {
		if err := connector.FinishVariant(vd.variant.ID, vd.variant.Status, vd.variant.Completed); err != nil {
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

	log.Infof("Job completed: ERRORED=%d, SUCCEEDED=%d, FAILED=%d\n", errorCount, successCount, failCount)

	var jobStatus lib.JobStatus
	switch {
	case errorCount > 0:
		jobStatus = lib.JOB_ERRORED
	case failCount > 0:
		jobStatus = lib.JOB_FAILED
	default:
		jobStatus = lib.JOB_SUCCESS

	}
	if err = connector.FinishJob(env[BazookaEnvJobID], jobStatus, time.Now()); err != nil {
		log.Fatal(err)
	}
	elapsed := time.Since(start)

	log.WithFields(log.Fields{
		"elapsed": elapsed,
	}).Info("Job Orchestration finished")
}

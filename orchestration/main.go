package main

import (
	"fmt"
	"log"
	"os"
	"time"

	lib "github.com/bazooka-ci/bazooka-lib"
	"github.com/bazooka-ci/bazooka-lib/mongo"
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

func main() {
	// TODO add validation
	start := time.Now()
	log.SetFlags(0)

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

	checkoutFolder := fmt.Sprintf(CheckoutFolderPattern, env[BazookaEnvHome])
	metaFolder := fmt.Sprintf(MetaFolderPattern, env[BazookaEnvHome])
	f := &SCMFetcher{
		Options: &FetchOptions{
			Scm:         env[BazookaEnvSCM],
			URL:         env[BazookaEnvSCMUrl],
			Reference:   env[BazookaEnvSCMReference],
			LocalFolder: checkoutFolder,
			MetaFolder:  metaFolder,
			KeyFile:     env[BazookaEnvSCMKeyfile],
			Env:         env,
		},
	}
	if err := f.Fetch(); err != nil {
		connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		log.Fatal(err)
	}

	p := &Parser{
		Options: &ParseOptions{
			InputFolder:    checkoutFolder,
			OutputFolder:   fmt.Sprintf(WorkdirFolderPattern, env[BazookaEnvHome]),
			DockerSock:     env[BazookaEnvDockerSock],
			HostBaseFolder: checkoutFolder,
			Env:            env,
		},
	}
	if err := p.Parse(); err != nil {
		connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		log.Fatal(err)
	}
	b := &Builder{
		Options: &BuildOptions{
			DockerfileFolder: fmt.Sprintf(WorkdirFolderPattern, BazookaInput),
			SourceFolder:     fmt.Sprintf(CheckoutFolderPattern, BazookaInput),
			ProjectID:        env[BazookaEnvProjectID],
			JobID:            env[BazookaEnvJobID],
		},
	}
	buildImages, err := b.Build()
	if err != nil {
		connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		log.Fatal(err)
	}

	r := &Runner{
		BuildImages: buildImages,
		Env:         env,
	}
	success, err := r.Run()
	if err != nil {
		connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		log.Fatal(err)
	}
	if success {
		connector.FinishJob(env[BazookaEnvJobID], lib.JOB_SUCCESS, time.Now())
	} else {
		connector.FinishJob(env[BazookaEnvJobID], lib.JOB_FAILED, time.Now())
	}

	elapsed := time.Since(start)
	log.Printf("Job Orchestration took %s", elapsed)

}

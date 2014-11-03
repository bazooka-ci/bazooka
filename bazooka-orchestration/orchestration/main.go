package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	CheckoutFolderPattern = "%s/source"
	WorkdirFolderPattern  = "%s/work"

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
)

func main() {
	// TODO add validation
	start := time.Now()

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
	f := &SCMFetcher{
		Options: &FetchOptions{
			Scm:         env[BazookaEnvSCM],
			URL:         env[BazookaEnvSCMUrl],
			Reference:   env[BazookaEnvSCMReference],
			LocalFolder: checkoutFolder,
			KeyFile:     env[BazookaEnvSCMKeyfile],
			Env:         env,
		},
	}
	if err := f.Fetch(); err != nil {
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
		log.Fatal(err)
	}
	b := &Builder{
		Options: &BuildOptions{
			DockerfileFolder: fmt.Sprintf(WorkdirFolderPattern, BazookaInput),
			SourceFolder:     fmt.Sprintf(CheckoutFolderPattern, BazookaInput),
			JobID:            env[BazookaEnvProjectID],
			VariantID:        env[BazookaEnvJobID],
		},
	}
	if err := b.Build(); err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	log.Printf("Job Orchestration took %s", elapsed)
}

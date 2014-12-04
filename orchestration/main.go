package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	docker "github.com/bywan/go-dockercommand"
	lib "github.com/haklop/bazooka/commons"
	l "github.com/haklop/bazooka/commons/logger"
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
		l.Init(ioutil.Discard, io.MultiWriter(os.Stdout, w), io.MultiWriter(os.Stdout, w), io.MultiWriter(os.Stdout, w), os.Stdout)
		connector.FeedLog(r, lib.LogEntry{
			ProjectID: env[BazookaEnvProjectID],
			JobID:     env[BazookaEnvJobID],
			Image:     "bazooka/orchestration",
		})
	}()

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
			l.Error.Fatal(mongoErr)
		}
		l.Error.Fatal(err)
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
	if err := p.Parse(containerLogger); err != nil {
		mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			l.Error.Fatal(mongoErr)
		}
		l.Error.Fatal(err)
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
		mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			l.Error.Fatal(mongoErr)
		}
		l.Error.Fatal(err)
	}

	r := &Runner{
		BuildImages: buildImages,
		Env:         env,
		Mongo:       connector,
	}
	success, err := r.Run(containerLogger)
	if err != nil {
		mongoErr := connector.FinishJob(env[BazookaEnvJobID], lib.JOB_ERRORED, time.Now())
		if mongoErr != nil {
			l.Error.Fatal(mongoErr)
		}
		l.Error.Fatal(err)
	}
	if success {
		err = connector.FinishJob(env[BazookaEnvJobID], lib.JOB_SUCCESS, time.Now())
	} else {
		err = connector.FinishJob(env[BazookaEnvJobID], lib.JOB_FAILED, time.Now())
	}
	if err != nil {
		l.Error.Fatal(err)
	}

	elapsed := time.Since(start)
	l.Info.Printf("Job Orchestration took %s", elapsed)

}

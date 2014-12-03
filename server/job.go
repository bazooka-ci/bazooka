package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	docker "github.com/bywan/go-dockercommand"
	lib "github.com/haklop/bazooka/commons"
	"github.com/haklop/bazooka/commons/mongo"
)

const (
	buildFolderPattern = "%s/build/%s/%s"     // $bzk_home/build/$projectId/$buildId
	logFolderPattern   = "%s/build/%s/%s/log" // $bzk_home/build/$projectId/$buildId/log
	// keyFolderPattern   = "%s/key/%s"         // $bzk_home/key/$keyName
)

func (c *context) startBuild(params map[string]string, body bodyFunc) (*response, error) {
	var startJob lib.StartJob

	body(&startJob)

	if len(startJob.ScmReference) == 0 {
		return badRequest("reference is mandatory")
	}

	project, err := c.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	client, err := docker.NewDocker(c.DockerEndpoint)
	if err != nil {
		return nil, err
	}

	runningJob := &lib.Job{
		ProjectID: project.ID,
		Started:   time.Now(),
	}

	if err := c.Connector.AddJob(runningJob); err != nil {
		return nil, err
	}

	buildFolder := fmt.Sprintf(buildFolderPattern, c.Env[BazookaEnvHome], runningJob.ProjectID, runningJob.ID)
	orchestrationEnv := map[string]string{
		"BZK_SCM":           project.ScmType,
		"BZK_SCM_URL":       project.ScmURI,
		"BZK_SCM_REFERENCE": startJob.ScmReference,
		"BZK_SCM_KEYFILE":   c.Env[BazookaEnvSCMKeyfile], //TODO use keyfile per project
		"BZK_HOME":          buildFolder,
		"BZK_PROJECT_ID":    project.ID,
		"BZK_JOB_ID":        runningJob.ID, // TODO handle job number and tasks and save it
		"BZK_DOCKERSOCK":    c.Env[BazookaEnvDockerSock],
		BazookaEnvMongoAddr: c.Env[BazookaEnvMongoAddr],
		BazookaEnvMongoPort: c.Env[BazookaEnvMongoPort],
	}

	container, err := client.Run(&docker.RunOptions{
		Image:       "bazooka/orchestration",
		VolumeBinds: []string{fmt.Sprintf("%s:/bazooka", buildFolder), fmt.Sprintf("%s:/var/run/docker.sock", c.Env[BazookaEnvDockerSock])},
		Env:         orchestrationEnv,
		Detach:      true,
	})

	logFolder := fmt.Sprintf(logFolderPattern, BazookaHome, runningJob.ProjectID, runningJob.ID)
	os.MkdirAll(logFolder, 0755)

	// Ensure directory exists
	err = os.MkdirAll(logFolder, 0755)
	if err != nil {
		return nil, err
	}
	logFileWriter, err := os.Create(logFolder + "/job.log")
	if err != nil {
		return nil, err
	}

	runningJob.OrchestrationID = container.ID()
	orchestrationLog := log.New(logFileWriter, "", log.LstdFlags)
	orchestrationLog.Printf("Start job %s on project %s with container %s\n", runningJob.ID, runningJob.ProjectID, runningJob.OrchestrationID)
	c.Connector.SetJobOrchestrationId(runningJob.ID, container.ID())
	if err != nil {
		orchestrationLog.Println(err.Error())
		return nil, err
	}

	r, w := io.Pipe()
	container.StreamLogs(w)
	go func(reader io.Reader, logFileWriter *os.File) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			orchestrationLog.Printf("%s \n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			orchestrationLog.Println("There was an error with the scanner in attached container", err)
		}
		logFileWriter.Close()
	}(r, logFileWriter)

	return accepted(runningJob, "/job/"+runningJob.ID)
}

func (c *context) getJob(params map[string]string, body bodyFunc) (*response, error) {

	job, err := c.Connector.GetJobByID(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("job not found")
	}

	return ok(&job)
}

func (c *context) getJobs(params map[string]string, body bodyFunc) (*response, error) {

	jobs, err := c.Connector.GetJobs(params["id"])
	if err != nil {
		return nil, err
	}

	return ok(&jobs)
}

func (c *context) getJobLog(params map[string]string, body bodyFunc) (*response, error) {

	log, err := c.Connector.GetLog(&mongo.LogExample{
		JobID: params["id"],
	})
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("log not found")
	}

	return ok(&log)
}

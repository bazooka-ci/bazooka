package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"bitbucket.org/bywan/bazooka-command/server/context"

	docker "github.com/bywan/go-dockercommand"
	"github.com/gorilla/mux"
	lib "github.com/haklop/bazooka/commons"
	"github.com/haklop/bazooka/commons/mongo"
)

const (
	buildFolderPattern = "%s/build/%s/%s"     // $bzk_home/build/$projectId/$buildId
	logFolderPattern   = "%s/build/%s/%s/log" // $bzk_home/build/$projectId/$buildId/log
	// keyFolderPattern   = "%s/key/%s"         // $bzk_home/key/$keyName
)

func (p *Context) startBuild(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var startJob lib.StartJob

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := decoder.Decode(&startJob)
	if err != nil {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "Invalid body : " + err.Error(),
		})
		return
	}

	if len(startJob.ScmReference) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&ErrorResponse{
			Code:    400,
			Message: "reference is mandatory",
		})

		return
	}

	project, err := p.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "project not found",
		})
		return
	}

	client, err := docker.NewDocker(p.DockerEndpoint)
	if err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	runningJob := &lib.Job{
		ProjectID: project.ID,
		Started:   time.Now(),
	}

	if err := p.Connector.AddJob(runningJob); err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	buildFolder := fmt.Sprintf(buildFolderPattern, p.Env[context.BazookaEnvHome], runningJob.ProjectID, runningJob.ID)
	orchestrationEnv := map[string]string{
		"BZK_SCM":                   "git",
		"BZK_SCM_URL":               project.ScmURI,
		"BZK_SCM_REFERENCE":         startJob.ScmReference,
		"BZK_SCM_KEYFILE":           p.Env[context.BazookaEnvSCMKeyfile], //TODO use keyfile per project
		"BZK_HOME":                  buildFolder,
		"BZK_PROJECT_ID":            project.ID,
		"BZK_JOB_ID":                runningJob.ID, // TODO handle job number and tasks and save it
		"BZK_DOCKERSOCK":            p.Env[context.BazookaEnvDockerSock],
		context.BazookaEnvMongoAddr: p.Env[context.BazookaEnvMongoAddr],
		context.BazookaEnvMongoPort: p.Env[context.BazookaEnvMongoPort],
	}

	container, err := client.Run(&docker.RunOptions{
		Image:       "bazooka/orchestration",
		VolumeBinds: []string{fmt.Sprintf("%s:/bazooka", buildFolder), fmt.Sprintf("%s:/var/run/docker.sock", p.Env[context.BazookaEnvDockerSock])},
		Env:         orchestrationEnv,
		Detach:      true,
	})

	logFolder := fmt.Sprintf(logFolderPattern, context.BazookaHome, runningJob.ProjectID, runningJob.ID)
	os.MkdirAll(logFolder, 0755)

	// Ensure directory exists
	err = os.MkdirAll(logFolder, 0755)
	if err != nil {
		log.Fatal(err)
	}
	logFileWriter, err := os.Create(logFolder + "/job.log")
	if err != nil {
		panic(err)
	}

	runningJob.OrchestrationID = container.ID()
	orchestrationLog := log.New(logFileWriter, "", log.LstdFlags)
	orchestrationLog.Printf("Start job %s on project %s with container %s\n", runningJob.ID, runningJob.ProjectID, runningJob.OrchestrationID)
	p.Connector.SetJobOrchestrationId(runningJob.ID, container.ID())
	if err != nil {
		orchestrationLog.Println(err.Error())
		context.WriteError(err, res, encoder)
		return
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

	res.Header().Set("Location", "/job/"+runningJob.ID)

	res.WriteHeader(202)
	encoder.Encode(runningJob)
}

func (p *Context) getJob(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	job, err := p.Connector.GetJobByID(params["job_id"])
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "job not found",
		})
		return
	}

	if params["id"] != job.ProjectID {
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "project not found",
		})
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&job)
}

func (p *Context) getJobs(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	jobs, err := p.Connector.GetJobs(params["id"])
	if err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&jobs)
}

func (p *Context) getJobLog(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	log, err := p.Connector.GetLog(&mongo.LogExample{
		JobID: params["job_id"],
	})
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&ErrorResponse{
			Code:    404,
			Message: "log not found",
		})
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&log)
}

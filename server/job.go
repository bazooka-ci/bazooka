package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/server/db"
)

func (c *context) startBitbucketJob(r *request) (*response, error) {
	var bitbucketPayload BitbucketPayload

	r.parseBody(&bitbucketPayload)

	if len(bitbucketPayload.Push.Changes) == 0 {
		return badRequest("no changes found in Bitbucket payload")
	}

	// Making the choice to return the result of the last change only
	// In most of the cases, there will be only 1 change
	var lastJobLaunchResponse *response
	var lastJobLaunchErr error
	for _, change := range bitbucketPayload.Push.Changes {
		changeType := change.New.Type
		if changeType == "annotated_tag" || changeType == "tag" {
			lastJobLaunchResponse, lastJobLaunchErr = c.startJob(r.vars, lib.StartJob{
				ScmReference: change.New.Name,
			}, "")
		} else if changeType == "branch" {
			lastJobLaunchResponse, lastJobLaunchErr = c.startJob(r.vars, lib.StartJob{
				ScmReference: change.New.Name,
			}, change.New.Target.Hash)
		} else {
			lastJobLaunchResponse, lastJobLaunchErr = badRequest(fmt.Sprintf("Change type %s unknow for Bitbucket payload", changeType))
		}
	}

	return lastJobLaunchResponse, lastJobLaunchErr

}

func (c *context) startStandardJob(r *request) (*response, error) {

	var startJob lib.StartJob

	r.parseBody(&startJob)

	if len(startJob.ScmReference) == 0 {
		return badRequest("reference is mandatory")
	}

	return c.startJob(r.vars, startJob, "")
}

func (c *context) startJob(params map[string]string, startJob lib.StartJob, commitID string) (*response, error) {

	project, err := c.connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	for _, v := range startJob.Parameters {
		if !strings.Contains(v, "=") {
			return nil, &errorResponse{400, fmt.Sprintf("Environment variable %v is empty", v)}
		}
	}

	runningJob := &lib.Job{
		Status:     lib.JOB_PENDING,
		ProjectID:  project.ID,
		Submitted:  time.Now(),
		Parameters: startJob.Parameters,
		SCMMetadata: lib.SCMMetadata{
			Reference: startJob.ScmReference,
			CommitID:  commitID,
		},
	}

	if err := c.connector.AddJob(runningJob); err != nil {
		return nil, fmt.Errorf("Failed to store job: %v", err)
	}

	jobBody, err := json.Marshal(runningJob)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal job: %v", err)
	}

	qid, err := c.queue.Put(jobBody, 0, 0*time.Second, 30*time.Second)

	if err != nil {
		return nil, fmt.Errorf("Failed to put job in queue: %v", err)
	}

	log.WithFields(log.Fields{
		"job_id":     runningJob.ID,
		"queue_id":   strconv.FormatUint(qid, 10),
		"project_id": runningJob.ProjectID,
	}).Info("Submitted job")

	return accepted(runningJob, "/job/"+runningJob.ID)
}

func (c *context) getJob(r *request) (*response, error) {

	job, err := c.connector.GetJobByID(r.vars["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("job not found")
	}

	return ok(&job)
}

func (c *context) getJobs(r *request) (*response, error) {

	jobs, err := c.connector.GetJobs(r.vars["id"])
	if err != nil {
		return nil, err
	}

	return ok(&jobs)
}

func (c *context) getAllJobs(r *request) (*response, error) {

	jobs, err := c.connector.GetAllJobs()
	if err != nil {
		return nil, err
	}

	return ok(&jobs)
}

func (c *context) getJobLog(r *request) (*response, error) {
	follow := len(r.query("follow")) > 0
	strictJson := len(r.query("strict-json")) > 0

	jid := r.vars["id"]

	job, err := c.connector.GetJobByID(jid)

	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("job not found")
	}

	w := r.w
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	logOutput := json.NewEncoder(w)

	query := &db.LogExample{
		JobID: jid,
	}

	logs, err := c.connector.GetLog(query)
	if !follow {
		logOutput.Encode(logs)
		return nil, nil
	}

	if strictJson {
		w.Write([]byte("["))
		defer w.Write([]byte("]"))
	}

	writtenEntries := 0
	for _, l := range logs {
		if writtenEntries > 0 && strictJson {
			w.Write([]byte(","))
		}
		logOutput.Encode(l)
		writtenEntries++
	}
	flushResponse(w)

	if job.Status != lib.JOB_RUNNING {
		return nil, nil
	}

	lastTime := jobLastLogTime(job, logs)

	for {
		time.Sleep(1000 * time.Millisecond)
		query.After = lastTime
		logs, err := c.connector.GetLog(query)
		if err != nil {
			log.Errorf("Error while retrieving logs: %v", err)
			return nil, nil
		}
		if len(logs) > 0 {
			lastTime = jobLastLogTime(job, logs)
			for _, l := range logs {
				if writtenEntries > 0 && strictJson {
					w.Write([]byte(","))
				}
				logOutput.Encode(l)
				writtenEntries++
			}
			flushResponse(w)
			continue
		}
		job, err := c.connector.GetJobByID(jid)
		if err != nil {
			log.Errorf("Error while retrieving job: %v", err)
			return nil, nil
		}
		if job.Status != lib.JOB_RUNNING {
			return nil, nil
		}
	}
}

func jobLastLogTime(job *lib.Job, logs []lib.LogEntry) time.Time {
	if len(logs) == 0 {
		return job.Started
	}
	return logs[len(logs)-1].Time
}

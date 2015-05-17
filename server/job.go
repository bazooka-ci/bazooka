package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/commons/mongo"
	docker "github.com/bywan/go-dockercommand"
)

const (
	buildFolderPattern = "%s/build/%s/%s"     // $bzk_home/build/$projectId/$buildId
	logFolderPattern   = "%s/build/%s/%s/log" // $bzk_home/build/$projectId/$buildId/log
)

func (c *context) startBitbucketJob(r *request) (*response, error) {
	var bitbucketPayload BitbucketPayload

	r.parseBody(&bitbucketPayload)

	if len(bitbucketPayload.Commits) == 0 {
		return badRequest("no commit found in Bitbucket payload")
	}

	//TODO(julienvey) Order by timestamp to find the last commit instead of trusting
	// Bitbucket to give us the commits in the right order

	if len(bitbucketPayload.Commits[0].RawNode) == 0 {
		return badRequest("RawNode is empty in Bitbucket payload")
	}

	return c.startJob(r.vars, lib.StartJob{
		ScmReference: bitbucketPayload.Commits[0].RawNode,
	}, "")

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

	orchestrationImage, err := c.Connector.GetImage("orchestration")
	if err != nil {
		return nil, &errorResponse{500, fmt.Sprintf("Failed to retrieve the orchestration image: %v", err)}
	}

	runningJob := &lib.Job{
		ProjectID:  project.ID,
		Started:    time.Now(),
		Parameters: startJob.Parameters,
		SCMMetadata: lib.SCMMetadata{
			Reference: startJob.ScmReference,
		},
	}

	if err := c.Connector.AddJob(runningJob); err != nil {
		return nil, err
	}

	var parametersAsBzkString []lib.BzkString
	for _, v := range startJob.Parameters {
		if !strings.Contains(v, "=") {
			return nil, &errorResponse{400, fmt.Sprintf("Environment variable %v is empty", v)}
		}
		name, value := lib.SplitNameValue(v)
		parametersAsBzkString = append(parametersAsBzkString, lib.BzkString{
			Name: name, Value: value, Secured: false,
		})
	}
	jobParameters, err := json.Marshal(parametersAsBzkString)
	if err != nil {
		return nil, err
	}

	var refToBuild string
	if len(commitID) > 0 {
		refToBuild = commitID
	} else {
		refToBuild = startJob.ScmReference
	}

	buildFolder := fmt.Sprintf(buildFolderPattern, c.Env[BazookaEnvHome], runningJob.ProjectID, runningJob.ID)
	orchestrationEnv := map[string]string{
		"BZK_SCM":            project.ScmType,
		"BZK_SCM_URL":        project.ScmURI,
		"BZK_SCM_REFERENCE":  refToBuild,
		"BZK_HOME":           buildFolder,
		"BZK_SRC":            buildFolder + "/source",
		"BZK_PROJECT_ID":     project.ID,
		"BZK_JOB_ID":         runningJob.ID,
		"BZK_DOCKERSOCK":     c.Env[BazookaEnvDockerSock],
		"BZK_JOB_PARAMETERS": string(jobParameters),
		BazookaEnvMongoAddr:  c.Env[BazookaEnvMongoAddr],
		BazookaEnvMongoPort:  c.Env[BazookaEnvMongoPort],
	}

	buildFolderLocal := fmt.Sprintf(buildFolderPattern, "/bazooka", runningJob.ProjectID, runningJob.ID)

	projectSSHKey, err := c.Connector.GetProjectKey(project.ID)
	if err != nil {
		_, keyNotFound := err.(*mongo.NotFoundError)
		if !keyNotFound {
			return nil, err
		}
		//Use Global Key if provided
		if len(c.Env[BazookaEnvSCMKeyfile]) > 0 {
			orchestrationEnv["BZK_SCM_KEYFILE"] = c.Env[BazookaEnvSCMKeyfile]
		}
	} else {
		err = os.MkdirAll(buildFolderLocal, 0644)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/key", buildFolderLocal), []byte(projectSSHKey.Content), 0600)
		if err != nil {
			return nil, err
		}
		orchestrationEnv["BZK_SCM_KEYFILE"] = fmt.Sprintf("%s/key", buildFolder)
	}

	projectCryptoKey, err := c.Connector.GetProjectCryptoKey(project.ID)

	if err != nil {
		_, keyNotFound := err.(*mongo.NotFoundError)
		if !keyNotFound {
			return nil, err
		}
	} else {
		err = os.MkdirAll(buildFolderLocal, 0644)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/crypto-key", buildFolderLocal), []byte(projectCryptoKey.Content), 0600)
		if err != nil {
			return nil, err
		}
		orchestrationEnv["BZK_CRYPTO_KEYFILE"] = fmt.Sprintf("%s/crypto-key", buildFolder)
	}

	orchestrationVolumes := []string{
		fmt.Sprintf("%s:/bazooka", buildFolder),
		fmt.Sprintf("%s:/var/run/docker.sock", c.Env[BazookaEnvDockerSock]),
	}

	reuseScmCheckout := project.Config["bzk.scm.reuse"] == "true"
	if reuseScmCheckout {
		hostSharedSourceFolder := fmt.Sprintf("%s/build/%s/source", c.Env[BazookaEnvHome], runningJob.ProjectID)
		containerSharedSourceFolder := fmt.Sprintf("/bazooka/build/%s/source", runningJob.ProjectID)

		_, err := os.Stat(containerSharedSourceFolder)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(containerSharedSourceFolder, 0644)
				if err != nil {
					return nil, fmt.Errorf("Failed to create a shared source directory for project %s, job %s: %v",
						runningJob.ProjectID, runningJob.ID, err)
				}
			} else {
				return nil, fmt.Errorf("Failed to stat the shared source directory for project %s, job %s: %v",
					runningJob.ProjectID, runningJob.ID, err)
			}
		}

		orchestrationEnv["BZK_SRC"] = hostSharedSourceFolder
		orchestrationEnv["BZK_REUSE_SCM_CHECKOUT"] = "1"

		orchestrationVolumes = append(orchestrationVolumes, fmt.Sprintf("%s:/bazooka/source", hostSharedSourceFolder))
	}

	container, err := client.Run(&docker.RunOptions{
		Image:       orchestrationImage,
		VolumeBinds: orchestrationVolumes,
		Env:         orchestrationEnv,
		Detach:      true,
	})

	// remove the container at the end of its execution
	go func(container *docker.Container) {
		exitCode, err := container.Wait()
		if err != nil {
			log.Errorf("Error while listening container %s", container.ID, err)
		}

		if exitCode != 0 {
			log.Errorf("Error during execution of Orchestrator container. Check Docker container logs, id is %s\n", container.ID())
			return
		}

		err = container.Remove(&docker.RemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil {
			log.Errorf("Cannot remove container %s", container.ID)
		}
	}(container)

	runningJob.OrchestrationID = container.ID()
	log.WithFields(log.Fields{
		"job_id":           runningJob.ID,
		"project_id":       runningJob.ProjectID,
		"orchestration_id": runningJob.OrchestrationID,
	}).Info("Starting job")

	err = c.Connector.SetJobOrchestrationId(runningJob.ID, container.ID())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return accepted(runningJob, "/job/"+runningJob.ID)
}

func (c *context) getJob(r *request) (*response, error) {

	job, err := c.Connector.GetJobByID(r.vars["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("job not found")
	}

	return ok(&job)
}

func (c *context) getJobs(r *request) (*response, error) {

	jobs, err := c.Connector.GetJobs(r.vars["id"])
	if err != nil {
		return nil, err
	}

	return ok(&jobs)
}

func (c *context) getAllJobs(r *request) (*response, error) {

	jobs, err := c.Connector.GetAllJobs()
	if err != nil {
		return nil, err
	}

	return ok(&jobs)
}

func (c *context) getJobLog(r *request) (*response, error) {
	follow := len(r.query("follow")) > 0
	strictJson := len(r.query("strict-json")) > 0

	jid := r.vars["id"]

	job, err := c.Connector.GetJobByID(jid)

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

	query := &mongo.LogExample{
		JobID: jid,
	}

	logs, err := c.Connector.GetLog(query)
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
		logs, err := c.Connector.GetLog(query)
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
		}
		job, err := c.Connector.GetJobByID(jid)
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

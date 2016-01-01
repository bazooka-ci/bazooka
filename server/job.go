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
	buildFolderPattern        = "%s/build/%s/%s"     // $bzk_home/build/$projectId/$buildId
	sharedSourceFolderPattern = "%s/build/%s/source" // $bzk_home/build/$projectId/source
	logFolderPattern          = "%s/build/%s/%s/log" // $bzk_home/build/$projectId/$buildId/log
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

	runningJob := &lib.Job{
		ProjectID:  project.ID,
		Started:    time.Now(),
		Parameters: startJob.Parameters,
		SCMMetadata: lib.SCMMetadata{
			Reference: startJob.ScmReference,
		},
	}
	if err := c.connector.AddJob(runningJob); err != nil {
		return nil, &errorResponse{500, fmt.Sprintf("Failed to add new job: %v", err)}
	}

	go c.runJob(&startJob, runningJob, commitID, project)

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

	query := &mongo.LogExample{
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

func (c *context) runJob(startJob *lib.StartJob, runningJob *lib.Job, commitID string, project *lib.Project) {
	client, err := docker.NewDocker(c.paths.dockerEndpoint.container)
	if err != nil {
		log.Errorf("Error creating new Docker client: %v", err)
	}

	orchestrationImage, err := c.connector.GetImage("orchestration")
	if err != nil {
		log.Errorf("Failed to retrieve the orchestration image: %v", err)
	}

	var parametersAsBzkString []lib.BzkString
	for _, v := range startJob.Parameters {
		if !strings.Contains(v, "=") {
			log.Errorf("Environment variable %v is empty", v)
		}
		name, value := lib.SplitNameValue(v)
		parametersAsBzkString = append(parametersAsBzkString, lib.BzkString{
			Name: name, Value: value, Secured: false,
		})
	}
	jobParameters, err := json.Marshal(parametersAsBzkString)
	if err != nil {
		log.Errorf("Error marshalling bzk job parameters: %v", err)
	}

	var refToBuild string
	if len(commitID) > 0 {
		refToBuild = commitID
	} else {
		refToBuild = startJob.ScmReference
	}

	buildFolder := path{
		host:      fmt.Sprintf(buildFolderPattern, c.paths.home.host, runningJob.ProjectID, runningJob.ID),
		container: fmt.Sprintf(buildFolderPattern, c.paths.home.container, runningJob.ProjectID, runningJob.ID),
	}
	orchestrationEnv := map[string]string{
		BazookaEnvApiUrl:     c.apiUrl,
		BazookaEnvSyslogUrl:  c.syslogUrl,
		BazookaEnvNetwork:    c.network,
		"BZK_SCM":            project.ScmType,
		"BZK_SCM_URL":        project.ScmURI,
		"BZK_SCM_REFERENCE":  refToBuild,
		"BZK_HOME":           buildFolder.host,
		"BZK_SRC":            buildFolder.host + "/source",
		"BZK_PROJECT_ID":     project.ID,
		"BZK_JOB_ID":         runningJob.ID,
		"BZK_DOCKERSOCK":     c.paths.dockerSock.host,
		"BZK_JOB_PARAMETERS": string(jobParameters),
		"BZK_FILE":           project.Config["bzk.file"],
	}

	projectSSHKey, err := c.connector.GetProjectKey(project.ID)
	if err != nil {
		_, keyNotFound := err.(*mongo.NotFoundError)
		if !keyNotFound {
			log.Errorf("Error getting Project SSH Key from Mongo: %v", err)
		}
		//Use Global Key if provided
		if len(c.paths.scmKey.host) > 0 {
			orchestrationEnv[BazookaEnvSCMKeyfile] = c.paths.scmKey.host
		}
	} else {
		err = os.MkdirAll(buildFolder.container, 0755)
		if err != nil {
			log.Errorf("Error creating build folder %s: %v", buildFolder.container, err)
		}

		keyFile := fmt.Sprintf("%s/key", buildFolder.container)
		err = ioutil.WriteFile(keyFile, []byte(projectSSHKey.Content), 0600)
		if err != nil {
			log.Errorf("Error writing key file in container %s: %v", keyFile, err)
		}
		orchestrationEnv[BazookaEnvSCMKeyfile] = fmt.Sprintf("%s/key", buildFolder.host)
	}

	projectCryptoKey, err := c.connector.GetProjectCryptoKey(project.ID)
	if err != nil {
		_, keyNotFound := err.(*mongo.NotFoundError)
		if !keyNotFound {
			log.Errorf("Error getting Project Crypto Key from Mongo: %v", err)
		}
	} else {
		err = os.MkdirAll(buildFolder.container, 0755)
		if err != nil {
			log.Errorf("Error creating build folder %s: %v", buildFolder.container, err)
		}

		cryptoFile := fmt.Sprintf("%s/crypto-key", buildFolder.container)
		err = ioutil.WriteFile(cryptoFile, []byte(projectCryptoKey.Content), 0600)
		if err != nil {
			log.Errorf("Error writing crypto file in container %s: %v", cryptoFile, err)
		}
		orchestrationEnv["BZK_CRYPTO_KEYFILE"] = fmt.Sprintf("%s/crypto-key", buildFolder.host)
	}

	orchestrationVolumes := []string{
		fmt.Sprintf("%s:/bazooka", buildFolder.host),
		fmt.Sprintf("%s:/var/run/docker.sock", c.paths.dockerSock.host),
	}

	reuseScmCheckout := project.Config["bzk.scm.reuse"] == "true"
	if reuseScmCheckout {
		sharedSourceFolder := path{
			host:      fmt.Sprintf(sharedSourceFolderPattern, c.paths.home.host, runningJob.ProjectID),
			container: fmt.Sprintf(sharedSourceFolderPattern, c.paths.home.container, runningJob.ProjectID),
		}

		_, err := os.Stat(sharedSourceFolder.container)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(sharedSourceFolder.container, 0755)
				if err != nil {
					log.Errorf("Failed to create a shared source directory for project %s, job %s: %v", runningJob.ProjectID, runningJob.ID, err)
				}
			} else {
				log.Errorf("Failed to stat the shared source directory for project %s, job %s: %v", runningJob.ProjectID, runningJob.ID, err)
			}
		}

		orchestrationEnv["BZK_SRC"] = sharedSourceFolder.host
		orchestrationEnv["BZK_REUSE_SCM_CHECKOUT"] = "1"

		orchestrationVolumes = append(orchestrationVolumes, fmt.Sprintf("%s:/bazooka/source", sharedSourceFolder.host))
	}

	container, err := client.Run(&docker.RunOptions{
		Image:         orchestrationImage.Image,
		VolumeBinds:   orchestrationVolumes,
		Env:           orchestrationEnv,
		Detach:        true,
		NetworkMode:   c.network,
		LoggingDriver: "syslog",
		LoggingDriverConfig: map[string]string{
			"syslog-address": c.syslogUrl,
			"syslog-tag": fmt.Sprintf("image=%s;project=%s;job=%s",
				orchestrationImage.Image, project.ID, runningJob.ID),
		},
	})

	if err != nil {
		log.Errorf("Failed to run the orchestration container for project %s, job %s: %v", runningJob.ProjectID, runningJob.ID, err)
	}

	defer lib.RemoveContainer(container)

	exitCode, err := container.Wait()
	if err != nil {
		log.Errorf("Error while waiting for container %s: %v", container.ID(), err)
	}

	if exitCode != 0 {
		log.Errorf("Error during execution of orchestration container with id %s, exit code is %d\n", container.ID(), exitCode)
	}

	runningJob.OrchestrationID = container.ID()
	log.WithFields(log.Fields{
		"job_id":           runningJob.ID,
		"project_id":       runningJob.ProjectID,
		"orchestration_id": runningJob.OrchestrationID,
	}).Info("Starting job")

	err = c.connector.SetJobOrchestrationId(runningJob.ID, container.ID())
	if err != nil {
		log.Error(err.Error())
	}
}

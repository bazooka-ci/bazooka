package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

func (c *context) startBitbucketJob(params map[string]string, body bodyFunc) (*response, error) {
	var bitbucketPayload BitbucketPayload

	body(&bitbucketPayload)

	if len(bitbucketPayload.Commits) == 0 {
		return badRequest("no commit found in Bitbucket payload")
	}

	//TODO(julienvey) Order by timestamp to find the last commit instead of trusting
	// Bitbucket to give us the commits in the right order

	if len(bitbucketPayload.Commits[0].RawNode) == 0 {
		return badRequest("RawNode is empty in Bitbucket payload")
	}

	return c.startJob(params, lib.StartJob{
		ScmReference: bitbucketPayload.Commits[0].RawNode,
	})

}

func (c *context) startGithubJob(params map[string]string, body bodyFunc) (*response, error) {
	var githubPayload GithubPayload

	body(&githubPayload)

	if len(githubPayload.HeadCommit.ID) == 0 {
		return badRequest("HeadCommit is empty in Github payload")
	}

	return c.startJob(params, lib.StartJob{
		ScmReference: githubPayload.HeadCommit.ID,
	})

}

func (c *context) startStandardJob(params map[string]string, body bodyFunc) (*response, error) {

	var startJob lib.StartJob

	body(&startJob)

	if len(startJob.ScmReference) == 0 {
		return badRequest("reference is mandatory")
	}

	return c.startJob(params, startJob)
}

func (c *context) startJob(params map[string]string, startJob lib.StartJob) (*response, error) {

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
		"BZK_HOME":          buildFolder,
		"BZK_PROJECT_ID":    project.ID,
		"BZK_JOB_ID":        runningJob.ID,
		"BZK_DOCKERSOCK":    c.Env[BazookaEnvDockerSock],
		BazookaEnvMongoAddr: c.Env[BazookaEnvMongoAddr],
		BazookaEnvMongoPort: c.Env[BazookaEnvMongoPort],
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
	fmt.Printf("Key: %s\n", projectCryptoKey)
	if err != nil {
		_, keyNotFound := err.(*mongo.NotFoundError)
		if !keyNotFound {
			return nil, err
		}
	} else {
		fmt.Printf("found\n")
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

	container, err := client.Run(&docker.RunOptions{
		Image: orchestrationImage,
		VolumeBinds: []string{
			fmt.Sprintf("%s:/bazooka", buildFolder),
			fmt.Sprintf("%s:/var/run/docker.sock", c.Env[BazookaEnvDockerSock]),
		},
		Env:    orchestrationEnv,
		Detach: true,
	})

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

func (c *context) getAllJobs(params map[string]string, body bodyFunc) (*response, error) {

	jobs, err := c.Connector.GetAllJobs()
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

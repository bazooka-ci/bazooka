package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bazooka-ci/bazooka/client"

	log "github.com/Sirupsen/logrus"
	docker "github.com/bywan/go-dockercommand"

	lib "github.com/bazooka-ci/bazooka/commons"
)

type reservedJob struct {
	id  uint64
	job *lib.Job
}

func (c *context) startJob(rj *reservedJob) {
	c.busy.Add(1)
	finished := make(chan bool)
	go c.jobHeartbeat(rj.id, finished)

	defer func() {
		finished <- true // to stop the periodic job touch
		if r := recover(); r != nil {
			log.Debugf("Recovered from: %v", r)
			switch r := r.(type) {
			case workerError:
				log.Errorf("Worker failed with %v", r.msg)

				if err := c.client.Internal.ResetJob(rj.job.ID); err != nil {
					log.Errorf("Error while resetting job %v: %v", rj.job.ID, err)
				}

				log.Infof("Releasing job %v", rj.id)
				c.releaseJob(rj.id)
				c.busy.Done()
			default:
				log.Errorf("Job %v errored with %v", rj.job.ID, r)

				log.Infof("Marking job %v as errored", rj.job.ID)
				err := c.client.Internal.MarkJobAsFinished(rj.job.ID, lib.JOB_ERRORED)
				if err != nil {
					log.Errorf("Failed to mark job %v as errored: %v", err)
					log.Infof("Releasing job %v", rj.id)
					c.releaseJob(rj.id)
					c.busy.Done()
					return
				}

				log.Infof("Deleting job %v from queue", rj.id)
				c.deleteJob(rj.id)
				c.busy.Done()
			}
		}
	}()

	job := rj.job
	log.WithFields(log.Fields{
		"ID":         job.ID,
		"QID":        fmt.Sprintf("%v", rj.id),
		"Params":     job.Parameters,
		"SCM_REF":    job.SCMMetadata.Reference,
		"SCM_COMMIT": job.SCMMetadata.CommitID,
		"Submitted":  job.Submitted,
	}).Infof("Starting job")

	if err := c.client.Internal.MarkJobAsStarted(job.ID); err != nil {
		workerPanic("Failed to mark the job as running: %v", err)
	}

	project, err := c.client.Project.Get(job.ProjectID)
	if err != nil {
		workerPanic("Failed to retrieve the job's project (id=%v): %v", job.ProjectID, err)
	}

	orchestrationImage, err := c.client.Image.Get("orchestration")
	if err != nil {
		workerPanic("Failed to retrieve the orchestration image: %v", err)
	}

	var parametersAsBzkString []lib.BzkString
	for _, v := range job.Parameters {
		name, value := lib.SplitNameValue(v)
		parametersAsBzkString = append(parametersAsBzkString, lib.BzkString{
			Name: name, Value: value, Secured: false,
		})
	}

	jobParameters, err := json.Marshal(parametersAsBzkString)
	if err != nil {
		jobPanic("Failed to serialize job parameters: %v", err)
	}

	refToBuild := job.SCMMetadata.CommitID
	if len(refToBuild) == 0 {
		refToBuild = job.SCMMetadata.Reference
	}

	buildFolder := c.buildFolder(job)

	orchestrationEnv := map[string]string{
		BazookaEnvServerApi:    c.serverApi,
		BazookaEnvServerSyslog: c.serverSyslog,
		BazookaEnvNetwork:      c.network,
		"BZK_SCM":              project.ScmType,
		"BZK_SCM_URL":          project.ScmURI,
		"BZK_SCM_REFERENCE":    refToBuild,
		"BZK_HOME":             buildFolder.host,
		"BZK_SRC":              buildFolder.host + "/source",
		"BZK_PROJECT_ID":       project.ID,
		"BZK_JOB_ID":           job.ID,
		"BZK_DOCKERSOCK":       c.paths.dockerSock.host,
		"BZK_JOB_PARAMETERS":   string(jobParameters),
	}

	projectSSHKey, err := c.client.Project.Key.Get(project.ID)
	if err != nil {
		if !client.IsNotFound(err) {
			workerPanic("Error retrieving the project SCM key: %v", err)
		}
		//Use Global Key if provided
		if len(c.paths.scmKey.host) > 0 {
			orchestrationEnv[BazookaEnvSCMKeyfile] = c.paths.scmKey.host
		}
	} else {
		err = os.MkdirAll(buildFolder.container, 0755)
		if err != nil {
			workerPanic("Error creating build folder %s: %v", buildFolder.container, err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/key", buildFolder.container), []byte(projectSSHKey.Content), 0600)
		if err != nil {
			workerPanic("Error writing the scm key: %v", err)
		}
		orchestrationEnv[BazookaEnvSCMKeyfile] = fmt.Sprintf("%s/key", buildFolder.host)
	}

	projectCryptoKey, err := c.client.Internal.GetProjectCryptoKey(project.ID)
	if err != nil {
		workerPanic("Error retrieving the project private key: %v", err)
	} else {
		err = os.MkdirAll(buildFolder.container, 0755)
		if err != nil {
			workerPanic("Error creating build folder %s: %v", buildFolder.container, err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/crypto-key", buildFolder.container), []byte(projectCryptoKey.Content), 0600)
		if err != nil {
			workerPanic("Error writing the crypto key: %v", err)
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
			host:      fmt.Sprintf(sharedSourceFolderPattern, c.paths.home.host, job.ProjectID),
			container: fmt.Sprintf(sharedSourceFolderPattern, c.paths.home.container, job.ProjectID),
		}

		_, err := os.Stat(sharedSourceFolder.container)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(sharedSourceFolder.container, 0755)
				if err != nil {
					workerPanic("Failed to create a shared source directory for project %s, job %s: %v",
						job.ProjectID, job.ID, err)
				}
			} else {
				workerPanic("Failed to stat the shared source directory for project %s, job %s: %v",
					job.ProjectID, job.ID, err)
			}
		}

		orchestrationEnv["BZK_SRC"] = sharedSourceFolder.host
		orchestrationEnv["BZK_REUSE_SCM_CHECKOUT"] = "1"

		orchestrationVolumes = append(orchestrationVolumes, fmt.Sprintf("%s:/bazooka/source", sharedSourceFolder.host))
	}

	log.WithFields(log.Fields{
		"Image":   orchestrationImage.Image,
		"Volumes": orchestrationVolumes,
		"Env":     orchestrationEnv,
	}).Info("Starting orchestration container")

	container, err := c.docker.Run(&docker.RunOptions{
		Image:         orchestrationImage.Image,
		VolumeBinds:   orchestrationVolumes,
		Env:           orchestrationEnv,
		Detach:        true,
		NetworkMode:   c.network,
		LoggingDriver: "syslog",
		LoggingDriverConfig: map[string]string{
			"syslog-address": c.serverSyslog,
			"syslog-tag": fmt.Sprintf("image=%s;project=%s;job=%s",
				orchestrationImage.Image, project.ID, job.ID),
		},
	})

	if err != nil {
		workerPanic("Error starting orchestration container: %v: %v", err, orchestrationVolumes)
	}

	defer lib.RemoveContainer(container)

	// remove the container at the end of its execution
	log.Infof("Waiting for job %v to run", job.ID)
	exitCode, err := container.Wait()
	if err != nil {
		workerPanic("Error while listening container %s", container.ID, err)
	}
	log.Infof("job %v finished with exit code %d", job.ID, exitCode)

	if exitCode != 0 {
		log.Errorf("Error during execution of orchestration container with id %s: exit code: %d", container.ID(), exitCode)
		jobPanic("orchestration exited with code %d", exitCode)
	}

	log.Infof("Deleting job %v from queue", rj.id)
	if err := c.deleteJob(rj.id); err != nil {
		log.Errorf("Cannot delete job %v: %v", rj.id, err)
	}
	log.Infof("Finished job %s", job.ID)
}

func (c *context) jobHeartbeat(qid uint64, finished chan bool) {
	ticker := time.Tick(10 * time.Second)
	for {
		select {
		case <-ticker:
			log.Debugf("Touching job %v", qid)
			c.touchJob(qid)
		case <-finished:
			log.Debugf("Stop touching %v", qid)
			return
		}
	}
}

type workerError struct {
	msg string
}

func workerPanic(format string, args ...interface{}) {
	panic(workerError{fmt.Sprintf(format, args...)})
}

func jobPanic(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

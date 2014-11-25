package project

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	docker "github.com/bywan/go-dockercommand"
	"github.com/gorilla/mux"
	lib "github.com/haklop/bazooka/commons"
	"github.com/haklop/bazooka/commons/mongo"
	"github.com/haklop/bazooka/server/context"
)

const (
	buildFolderPattern = "%s/build/%s/%s"     // $bzk_home/build/$projectId/$buildId
	logFolderPattern   = "%s/build/%s/%s/log" // $bzk_home/build/$projectId/$buildId/log
	// keyFolderPattern   = "%s/key/%s"         // $bzk_home/key/$keyName
)

type Handlers struct {
	mongoConnector *mongo.MongoConnector
	env            map[string]string
	dockerEndpoint string
}

func (p *Handlers) SetHandlers(r *mux.Router, serverContext context.Context) {
	p.mongoConnector = serverContext.Connector
	p.env = serverContext.Env
	p.dockerEndpoint = serverContext.DockerEndpoint

	r.HandleFunc("/", p.createProject).Methods("POST")
	r.HandleFunc("/", p.getProjects).Methods("GET")
	r.HandleFunc("/{id}", p.getProject).Methods("GET")
	r.HandleFunc("/{id}/job", p.startBuild).Methods("POST")
	r.HandleFunc("/{id}/job/", p.getJobs).Methods("GET")
	r.HandleFunc("/{id}/job/{job_id}", p.getJob).Methods("GET")
	r.HandleFunc("/{id}/job/{job_id}/log", p.getJobLog).Methods("GET")
	r.HandleFunc("/{id}/job/{job_id}/variant", p.getVariants).Methods("GET")
	r.HandleFunc("/{id}/job/{job_id}/variant/{variant_id}", p.getVariant).Methods("GET")
}

func (p *Handlers) createProject(res http.ResponseWriter, req *http.Request) {
	var project lib.Project

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := decoder.Decode(&project)

	if err != nil {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "Unable to decode your json : " + err.Error(),
		})
		return
	}

	if len(project.ScmURI) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "scm_uri is mandatory",
		})

		return
	}

	if len(project.ScmType) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "scm_type is mandatory",
		})

		return
	}

	existantProject, err := p.mongoConnector.GetProject(project.ScmType, project.ScmURI)
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
	}

	if len(existantProject.ScmURI) > 0 {
		res.WriteHeader(409)
		encoder.Encode(&context.ErrorResponse{
			Code:    409,
			Message: "scm_uri is already known",
		})

		return
	}

	// TODO : validate scm_type
	// TODO : validate data by scm_type

	err = p.mongoConnector.AddProject(&project)
	res.Header().Set("Location", "/project/"+project.ID)

	res.WriteHeader(201)
	encoder.Encode(&project)

}

func (p *Handlers) getProject(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	project, err := p.mongoConnector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&context.ErrorResponse{
			Code:    404,
			Message: "project not found",
		})

		return
	}

	res.WriteHeader(200)
	encoder.Encode(&project)
}

func (p *Handlers) getProjects(res http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	projects, err := p.mongoConnector.GetProjects()
	if err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&projects)
}

func (p *Handlers) startBuild(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var startJob lib.StartJob

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := decoder.Decode(&startJob)
	if err != nil {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "Invalid body : " + err.Error(),
		})
		return
	}

	if len(startJob.ScmReference) == 0 {
		res.WriteHeader(400)
		encoder.Encode(&context.ErrorResponse{
			Code:    400,
			Message: "reference is mandatory",
		})

		return
	}

	project, err := p.mongoConnector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&context.ErrorResponse{
			Code:    404,
			Message: "project not found",
		})
		return
	}

	client, err := docker.NewDocker(p.dockerEndpoint)
	if err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	runningJob := &lib.Job{
		ID:        strconv.FormatInt(time.Now().Unix(), 10),
		ProjectID: project.ID,
		Started:   time.Now(),
	}

	if err := p.mongoConnector.AddJob(runningJob); err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	buildFolder := fmt.Sprintf(buildFolderPattern, p.env[context.BazookaEnvHome], runningJob.ProjectID, runningJob.ID)
	orchestrationEnv := map[string]string{
		"BZK_SCM":                   "git",
		"BZK_SCM_URL":               project.ScmURI,
		"BZK_SCM_REFERENCE":         startJob.ScmReference,
		"BZK_SCM_KEYFILE":           p.env[context.BazookaEnvSCMKeyfile], //TODO use keyfile per project
		"BZK_HOME":                  buildFolder,
		"BZK_PROJECT_ID":            project.ID,
		"BZK_JOB_ID":                runningJob.ID, // TODO handle job number and tasks and save it
		"BZK_DOCKERSOCK":            p.env[context.BazookaEnvDockerSock],
		context.BazookaEnvMongoAddr: p.env[context.BazookaEnvMongoAddr],
		context.BazookaEnvMongoPort: p.env[context.BazookaEnvMongoPort],
	}

	container, err := client.Run(&docker.RunOptions{
		Image:       "bazooka/orchestration",
		VolumeBinds: []string{fmt.Sprintf("%s:/bazooka", buildFolder), fmt.Sprintf("%s:/var/run/docker.sock", p.env[context.BazookaEnvDockerSock])},
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
	p.mongoConnector.SetJobOrchestrationId(runningJob.ID, container.ID())
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

	res.Header().Set("Location", "/project/"+project.ID+"/job/"+runningJob.ID)

	res.WriteHeader(202)
	encoder.Encode(runningJob)
}

func (p *Handlers) getJob(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	job, err := p.mongoConnector.GetJobByID(params["job_id"])
	if err != nil {
		if err.Error() != "not found" {
			context.WriteError(err, res, encoder)
			return
		}
		res.WriteHeader(404)
		encoder.Encode(&context.ErrorResponse{
			Code:    404,
			Message: "job not found",
		})
		return
	}

	if params["id"] != job.ProjectID {
		res.WriteHeader(404)
		encoder.Encode(&context.ErrorResponse{
			Code:    404,
			Message: "project not found",
		})
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&job)
}

func (p *Handlers) getJobs(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	encoder := json.NewEncoder(res)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	jobs, err := p.mongoConnector.GetJobs(params["id"])
	if err != nil {
		context.WriteError(err, res, encoder)
		return
	}

	res.WriteHeader(200)
	encoder.Encode(&jobs)
}

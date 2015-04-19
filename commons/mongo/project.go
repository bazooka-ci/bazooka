package mongo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	lib "github.com/bazooka-ci/bazooka/commons"
	"gopkg.in/mgo.v2/bson"
)

func (c *MongoConnector) HasProject(name string, scmType string, scmURI string) (bool, error) {
	request := bson.M{}

	if len(name) > 0 {
		request["name"] = name
	}
	if len(scmType) > 0 {
		request["scm_type"] = scmType
	}
	if len(scmURI) > 0 {
		request["scm_uri"] = scmURI
	}

	count, err := c.database.C("projects").Find(request).Count()
	return count > 0, err
}

func (c *MongoConnector) GetProjectById(id string) (*lib.Project, error) {
	result := &lib.Project{}
	if err := c.ByIdOrName("projects", id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) GetProjects() ([]*lib.Project, error) {
	result := []*lib.Project{}

	err := c.database.C("projects").Find(bson.M{}).All(&result)
	return result, err
}

func (c *MongoConnector) AddProject(project *lib.Project) error {
	var err error
	if project.ID, err = c.randomId(); err != nil {
		return err
	}

	return c.database.C("projects").Insert(project)
}

func (c *MongoConnector) SetProjectConfig(id string, config map[string]string) error {
	proj, err := c.GetProjectById(id)
	if err != nil {
		return err
	}
	selector := bson.M{
		"id": proj.ID,
	}
	request := bson.M{
		"$set": bson.M{
			"config": config,
		},
	}
	return c.database.C("projects").Update(selector, request)
}

func (c *MongoConnector) SetProjectConfigKey(id, key, value string) error {
	proj, err := c.GetProjectById(id)
	if err != nil {
		return err
	}
	selector := bson.M{
		"id": proj.ID,
	}
	request := bson.M{
		"$set": bson.M{
			fmt.Sprintf("config.%s", key): value,
		},
	}
	return c.database.C("projects").Update(selector, request)
}

func (c *MongoConnector) UnsetProjectConfigKey(id, key string) error {
	proj, err := c.GetProjectById(id)
	if err != nil {
		return err
	}
	selector := bson.M{
		"id": proj.ID,
	}
	request := bson.M{
		"$unset": bson.M{
			fmt.Sprintf("config.%s", key): "",
		},
	}
	return c.database.C("projects").Update(selector, request)
}

func (c *MongoConnector) AddJob(job *lib.Job) error {
	var err error
	if job.ID, err = c.randomId(); err != nil {
		return err
	}

	if len(job.Status) == 0 {
		job.Status = lib.JOB_RUNNING
	}

	query := bson.M{"id": job.ProjectID}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"job_counter": 1}},
		ReturnNew: true,
	}
	var proj lib.Project
	_, err = c.database.C("projects").Find(query).Apply(change, &proj)
	if err != nil {
		return fmt.Errorf("Error generating the job number: %v", err)
	}
	job.Number = proj.JobCounter

	return c.database.C("jobs").Insert(job)
}

func (c *MongoConnector) AddVariant(variant *lib.Variant) error {
	var err error
	if variant.ID, err = c.randomId(); err != nil {
		return err
	}
	if len(variant.Status) == 0 {
		variant.Status = lib.JOB_RUNNING
	}
	return c.database.C("variants").Insert(variant)
}

func (c *MongoConnector) AddLog(log *lib.LogEntry) error {
	var err error
	if log.ID, err = c.randomId(); err != nil {
		return err
	}
	return c.database.C("logs").Insert(log)
}

type LogExample struct {
	ProjectID string
	JobID     string
	VariantID string
	Images    []string
}

func (c *MongoConnector) GetLog(like *LogExample) ([]lib.LogEntry, error) {
	result := []lib.LogEntry{}
	request := bson.M{}
	if len(like.ProjectID) > 0 {
		proj, err := c.GetProjectById(like.ProjectID)
		if err != nil {
			return nil, err
		}

		request["project_id"] = proj.ID
	}
	if len(like.JobID) > 0 {
		job, err := c.GetJobByID(like.JobID)
		if err != nil {
			return nil, err
		}
		request["job_id"] = job.ID
	}
	if len(like.VariantID) > 0 {
		v, err := c.GetVariantByID(like.VariantID)
		if err != nil {
			return nil, err
		}
		request["variant_id"] = v.ID
	}

	if len(like.Images) > 0 {
		request["image"] = bson.M{
			"$in": like.Images,
		}
	}
	err := c.database.C("logs").Find(request).All(&result)
	return result, err
}

func (c *MongoConnector) FeedLog(r io.Reader, template lib.LogEntry) {
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			message := scanner.Text()
			thisTemplate := template

			regLogLevel, _ := regexp.Compile(`^\s*\[(\S+)\].*$`)  // Eg. [INFO] My message
			regMeta, _ := regexp.Compile(`^\s*\<(\S+):(.*)>\s*$`) // Eg. <CMD:go test -v ./...>

			switch {
			case regLogLevel.MatchString(message):
				submatchs := regLogLevel.FindStringSubmatch(message)
				logLevel := submatchs[len(submatchs)-1]
				thisTemplate.Level = logLevel
				thisTemplate.Message = strings.TrimSpace(message[len(logLevel)+2:])
			case regMeta.MatchString(message):
				submatchs := regMeta.FindStringSubmatch(message)
				instructionType := submatchs[1]
				instructionValue := submatchs[2]
				switch instructionType {
				case "CMD":
					thisTemplate.Command = instructionValue
				case "PHASE":
					thisTemplate.Phase = instructionValue
				default:
					thisTemplate.Message = strings.TrimSpace(message)
				}
			default:
				thisTemplate.Message = strings.TrimSpace(message)
			}

			thisTemplate.Time = time.Now()
			c.AddLog(&thisTemplate)
		}
		if err := scanner.Err(); err != nil {
			log.Println("There was an error with the scanner", err)
		}
	}(r)
}

func (c *MongoConnector) SetJobOrchestrationId(id string, orchestrationId string) error {
	selector := bson.M{
		"id": id,
	}
	request := bson.M{
		"$set": bson.M{"orchestration_id": orchestrationId},
	}
	err := c.database.C("jobs").Update(selector, request)

	return err
}

func (c *MongoConnector) FinishJob(id string, status lib.JobStatus, completed time.Time) error {
	request := bson.M{
		"$set": bson.M{
			"status":    status,
			"completed": completed,
		},
	}
	err := c.database.C("jobs").Update(c.idLike(id), request)

	return err
}

func (c *MongoConnector) AddJobSCMMetadata(id string, metadata *lib.SCMMetadata) error {
	selector := bson.M{
		"id": id,
	}
	request := bson.M{
		"$set": bson.M{
			"scm_metadata": metadata,
		},
	}
	err := c.database.C("jobs").Update(selector, request)

	return err
}

func (c *MongoConnector) FinishVariant(id string, status lib.JobStatus, completed time.Time) error {
	request := bson.M{
		"$set": bson.M{
			"status":    status,
			"completed": completed,
		},
	}
	err := c.database.C("variants").Update(c.idLike(id), request)

	return err
}

func (c *MongoConnector) GetJobByID(id string) (*lib.Job, error) {
	result := &lib.Job{}
	if err := c.ById("jobs", id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) GetVariantByID(id string) (*lib.Variant, error) {
	result := &lib.Variant{}
	if err := c.ById("variants", id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) GetJobs(projectID string) ([]*lib.Job, error) {
	proj, err := c.GetProjectById(projectID)
	if err != nil {
		return nil, err
	}

	result := []*lib.Job{}
	err = c.database.C("jobs").Find(bson.M{
		"project_id": proj.ID,
	}).All(&result)
	return result, err
}

func (c *MongoConnector) GetAllJobs() ([]*lib.Job, error) {
	result := []*lib.Job{}
	err := c.database.C("jobs").Find(bson.M{}).All(&result)
	return result, err
}

func (c *MongoConnector) GetVariants(jobID string) ([]*lib.Variant, error) {
	job, err := c.GetJobByID(jobID)
	if err != nil {
		return nil, err
	}

	result := []*lib.Variant{}

	err = c.database.C("variants").Find(bson.M{
		"job_id": job.ID,
	}).All(&result)
	return result, err
}

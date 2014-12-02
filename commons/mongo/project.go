package mongo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"

	lib "github.com/haklop/bazooka/commons"
	"gopkg.in/mgo.v2/bson"
)

func (c *MongoConnector) GetProject(scmType string, scmURI string) (lib.Project, error) {
	result := lib.Project{}

	request := bson.M{
		"scm_uri":  scmURI,
		"scm_type": scmType,
	}
	err := c.database.C("projects").Find(request).One(&result)
	fmt.Printf("retrieve project: %#v\n", result)
	return result, err
}

func (c *MongoConnector) GetProjectById(id string) (*lib.Project, error) {
	result := &lib.Project{}
	if err := c.ById("projects", id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) GetProjects() ([]*lib.Project, error) {
	result := []*lib.Project{}

	err := c.database.C("projects").Find(bson.M{}).All(&result)
	fmt.Printf("retrieve projects: %#v\n", result)
	return result, err
}

func (c *MongoConnector) AddProject(project *lib.Project) error {
	var err error
	if project.ID, err = c.randomId(); err != nil {
		return err
	}

	fmt.Printf("add project: %#v\n", project)
	return c.database.C("projects").Insert(project)
}

func (c *MongoConnector) AddJob(job *lib.Job) error {
	var err error
	if job.ID, err = c.randomId(); err != nil {
		return err
	}

	if len(job.Status) == 0 {
		job.Status = lib.JOB_RUNNING
	}
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
			template.Message = scanner.Text()
			template.Time = time.Now()
			c.AddLog(&template)
		}
		if err := scanner.Err(); err != nil {
			log.Println("There was an error with the scanner in attached container", err)
		}
	}(r)
}

func (c *MongoConnector) SetJobOrchestrationId(id string, orchestrationId string) error {
	fmt.Printf("set job: %v orchestration id to %v", id, orchestrationId)
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
	fmt.Printf("finish job: %v with status %v\n", id, status)
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
	fmt.Printf("adding metadata: %v for job %v\n", metadata, id)
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
	fmt.Printf("finish variant: %v with status %v\n", id, status)
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
	fmt.Printf("retrieve jobs: %#v\n", result)
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
	fmt.Printf("retrieve variants: %#v\n", result)
	return result, err
}

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
	fmt.Printf("retrieve project: %#v", result)
	return result, err
}

func (c *MongoConnector) GetProjectById(id string) (lib.Project, error) {
	result := lib.Project{}

	request := bson.M{
		"id": id,
	}
	err := c.database.C("projects").Find(request).One(&result)
	fmt.Printf("retrieve project: %#v", result)
	return result, err
}

func (c *MongoConnector) GetProjects() ([]lib.Project, error) {
	result := []lib.Project{}

	err := c.database.C("projects").Find(bson.M{}).All(&result)
	fmt.Printf("retrieve projects: %#v", result)
	return result, err
}

func (c *MongoConnector) AddProject(project *lib.Project) error {
	i := bson.NewObjectId()
	project.ID = i.Hex()

	fmt.Printf("add project: %#v", project)
	err := c.database.C("projects").Insert(project)

	return err
}

func (c *MongoConnector) AddJob(job *lib.Job) error {
	fmt.Printf("add job: %#v", job)
	if len(job.Status) == 0 {
		job.Status = lib.JOB_RUNNING
	}
	err := c.database.C("jobs").Insert(job)

	return err
}

func (c *MongoConnector) AddVariant(variant *lib.Variant) error {
	i := bson.NewObjectId()
	variant.ID = i.Hex()

	fmt.Printf("add variant: %#v", variant)
	if len(variant.Status) == 0 {
		variant.Status = lib.JOB_RUNNING
	}
	err := c.database.C("variants").Insert(variant)

	return err
}

func (c *MongoConnector) AddLog(log *lib.LogEntry) error {
	i := bson.NewObjectId()
	log.ID = i.Hex()

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
		request["project_id"] = like.ProjectID
	}
	if len(like.JobID) > 0 {
		request["job_id"] = like.JobID
	}
	if len(like.VariantID) > 0 {
		request["variant_id"] = like.VariantID
	}

	if len(like.Images) > 0 {
		request["image"] = bson.M{
			"$in": like.Images,
		}
	}
	err := c.database.C("logs").Find(request).All(&result)
	fmt.Printf("retrieve projects: %#v", result)
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
	selector := bson.M{
		"id": id,
	}
	request := bson.M{
		"$set": bson.M{
			"status":    status,
			"completed": completed,
		},
	}
	err := c.database.C("jobs").Update(selector, request)

	return err
}

func (c *MongoConnector) FinishVariant(id string, status lib.JobStatus, completed time.Time) error {
	fmt.Printf("finish variant: %v with status %v\n", id, status)
	selector := bson.M{
		"id": id,
	}
	request := bson.M{
		"$set": bson.M{
			"status":    status,
			"completed": completed,
		},
	}
	err := c.database.C("variants").Update(selector, request)

	return err
}

func (c *MongoConnector) GetJobByID(id string) (lib.Job, error) {
	result := lib.Job{}

	request := bson.M{
		"id": id,
	}
	err := c.database.C("jobs").Find(request).One(&result)
	fmt.Printf("retrieve job: %#v", result)
	return result, err
}

func (c *MongoConnector) GetVariantByID(id string) (lib.Variant, error) {
	result := lib.Variant{}

	request := bson.M{
		"id": id,
	}
	err := c.database.C("variants").Find(request).One(&result)
	fmt.Printf("retrieve variant: %#v", result)
	return result, err
}

func (c *MongoConnector) GetJobs(projectID string) ([]lib.Job, error) {
	result := []lib.Job{}

	err := c.database.C("jobs").Find(bson.M{
		"project_id": projectID,
	}).All(&result)
	fmt.Printf("retrieve jobs: %#v", result)
	return result, err
}

func (c *MongoConnector) GetVariants(jobID string) ([]lib.Variant, error) {
	result := []lib.Variant{}

	err := c.database.C("variants").Find(bson.M{
		"job_id": jobID,
	}).All(&result)
	fmt.Printf("retrieve variants: %#v", result)
	return result, err
}

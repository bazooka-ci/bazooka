package db

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	lib "github.com/bazooka-ci/bazooka/commons"
	"gopkg.in/mgo.v2/bson"
)

func (c *MongoConnector) HasProject(name string) (bool, error) {
	request := bson.M{}

	if len(name) > 0 {
		request["name"] = name
	}

	count, err := c.database.C("projects").Find(request).Count()
	return count > 0, err
}

func (c *MongoConnector) GetProjectById(id string) (*lib.Project, error) {
	result := &lib.Project{}
	if err := c.selectOneByIdOrName("projects", id, result); err != nil {
		return nil, err
	}
	result.Config = unescapeDotsInMap(result.Config)

	return result, nil
}

func (c *MongoConnector) GetProjects() ([]*lib.Project, error) {
	result := []*lib.Project{}

	err := c.database.C("projects").Find(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	for _, project := range result {
		project.Config = unescapeDotsInMap(project.Config)
	}
	return result, nil
}

func (c *MongoConnector) GetProjectsWithStatus() ([]*lib.ProjectWithStatus, error) {
	projects, err := c.GetProjects()
	if err != nil {
		return nil, err
	}

	jobs := []*lib.Job{}
	if err := c.database.C("jobs").
		Pipe([]bson.M{
		{"$sort": bson.M{"started": -1}},
		{
			"$group": bson.M{
				"_id":        "$project_id",
				"project_id": bson.M{"$first": "$project_id"},
				"started":    bson.M{"$first": "$started"},
				"completed":  bson.M{"$first": "$completed"},
				"number":     bson.M{"$first": "$number"},
				"status":     bson.M{"$first": "$status"},
			},
		},
	}).All(&jobs); err != nil {
		return nil, err
	}

	indexed := map[string]*lib.Job{}
	for _, job := range jobs {
		indexed[job.ProjectID] = job
	}

	result := []*lib.ProjectWithStatus{}
	for _, project := range projects {
		s := &lib.ProjectWithStatus{
			Project: project,
		}
		if job, found := indexed[project.ID]; found {
			s.LastJob = job
		}
		result = append(result, s)
	}

	return result, nil
}

func (c *MongoConnector) AddProject(project *lib.Project) error {
	var err error
	if project.ID, err = c.randomId(); err != nil {
		return err
	}

	if project.HookKey, err = c.randomId(); err != nil {
		return err
	}

	project.Config = escapeDotsInMap(project.Config)
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
	escapedConfig := map[string]string{}
	for k, v := range config {
		escapedConfig[escapeDots(k)] = v
	}
	request := bson.M{
		"$set": bson.M{
			"config": escapedConfig,
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
			fmt.Sprintf("config.%s", escapeDots(key)): value,
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
			fmt.Sprintf("config.%s", escapeDots(key)): "",
		},
	}
	return c.database.C("projects").Update(selector, request)
}

const (
	escapedDot = "//"
)

func escapeDots(s string) string {
	return strings.Replace(s, ".", escapedDot, -1)
}

func escapeDotsInMap(m map[string]string) map[string]string {
	u := map[string]string{}
	for k, v := range m {
		u[escapeDots(k)] = v
	}
	return u
}

func unescapeDots(s string) string {
	return strings.Replace(s, escapedDot, ".", -1)
}

func unescapeDotsInMap(m map[string]string) map[string]string {
	u := map[string]string{}
	for k, v := range m {
		u[unescapeDots(k)] = v
	}
	return u
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
	After     time.Time
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

	if !like.After.IsZero() {
		request["time"] = bson.M{
			"$gt": like.After,
		}
	}

	err := c.database.C("logs").Find(request).All(&result)
	return result, err
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

func (c *MongoConnector) MarkJobAsStarted(id string, started time.Time) error {
	request := bson.M{
		"$set": bson.M{
			"status":  lib.JOB_RUNNING,
			"started": started,
		},
	}
	return c.database.C("jobs").Update(c.fieldStartsWith("id", id), request)
}

func (c *MongoConnector) MarkJobAsFinished(id string, status lib.JobStatus, completed time.Time) error {
	request := bson.M{
		"$set": bson.M{
			"status":    status,
			"completed": completed,
		},
	}
	return c.database.C("jobs").Update(c.fieldStartsWith("id", id), request)
}

func (c *MongoConnector) ResetJob(id string) error {
	request := bson.M{
		"$set": bson.M{
			"status": lib.JOB_PENDING,
		},
		"$unset": bson.M{
			"started":      "",
			"completed":    "",
			"scm_metadata": "",
		},
	}
	return c.database.C("jobs").Update(c.fieldStartsWith("id", id), request)
}

func (c *MongoConnector) AddJobSCMMetadata(id string, metadata *lib.SCMMetadata) error {
	job, err := c.GetJobByID(id)
	if err != nil {
		return err
	}

	// Do not override SCMMetadata reference if present in database
	if len(job.SCMMetadata.Reference) > 0 {
		metadata.Reference = job.SCMMetadata.Reference
	}

	selector := bson.M{
		"id": id,
	}
	request := bson.M{
		"$set": bson.M{
			"scm_metadata": metadata,
		},
	}
	return c.database.C("jobs").Update(selector, request)
}

func (c *MongoConnector) FinishVariant(id string, status lib.JobStatus, completed time.Time, artifacts []string) error {
	request := bson.M{
		"$set": bson.M{
			"status":    status,
			"completed": completed,
			"artifacts": artifacts,
		},
	}
	return c.database.C("variants").Update(c.fieldStartsWith("id", id), request)
}

func (c *MongoConnector) GetJobByID(id string) (*lib.Job, error) {
	result := &lib.Job{}
	if err := c.selectOneByFieldLike("jobs", "id", id, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MongoConnector) GetVariantByID(id string) (*lib.Variant, error) {
	result := &lib.Variant{}
	if err := c.selectOneByFieldLike("variants", "id", id, result); err != nil {
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

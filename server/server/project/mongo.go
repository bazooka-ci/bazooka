package project

import (
	"fmt"

	lib "github.com/bazooka-ci/bazooka-lib"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoConnector struct {
	Database *mgo.Database
}

func (c *mongoConnector) GetProject(scmType string, scmURI string) (lib.Project, error) {
	result := lib.Project{}

	request := bson.M{
		"scm_uri":  scmURI,
		"scm_type": scmType,
	}
	err := c.Database.C("projects").Find(request).One(&result)
	fmt.Printf("retrieve project: %#v", result)
	return result, err
}

func (c *mongoConnector) GetProjectById(id string) (lib.Project, error) {
	result := lib.Project{}

	request := bson.M{
		"id": id,
	}
	err := c.Database.C("projects").Find(request).One(&result)
	fmt.Printf("retrieve project: %#v", result)
	return result, err
}

func (c *mongoConnector) GetProjects() ([]lib.Project, error) {
	result := []lib.Project{}

	err := c.Database.C("projects").Find(bson.M{}).All(&result)
	fmt.Printf("retrieve projects: %#v", result)
	return result, err
}

func (c *mongoConnector) AddProject(project *lib.Project) error {
	i := bson.NewObjectId()
	project.ID = i.Hex()

	fmt.Printf("add project: %#v", project)
	err := c.Database.C("projects").Insert(project)

	return err
}

func (c *mongoConnector) AddJob(job *lib.Job) error {
	fmt.Printf("add job: %#v", job)
	err := c.Database.C("jobs").Insert(job)

	return err
}

func (c *mongoConnector) GetJobByID(id string) (lib.Job, error) {
	result := lib.Job{}

	request := bson.M{
		"id": id,
	}
	err := c.Database.C("jobs").Find(request).One(&result)
	fmt.Printf("retrieve job: %#v", result)
	return result, err
}

func (c *mongoConnector) GetJobs() ([]lib.Job, error) {
	result := []lib.Job{}

	err := c.Database.C("jobs").Find(bson.M{}).All(&result)
	fmt.Printf("retrieve jobs: %#v", result)
	return result, err
}

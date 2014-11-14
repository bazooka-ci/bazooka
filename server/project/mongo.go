package project

import (
	"fmt"

	bazooka "github.com/haklop/bazooka/commons"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoConnector struct {
	Database *mgo.Database
}

func (c *mongoConnector) GetProject(scmType string, scmURI string) (bazooka.Project, error) {
	result := bazooka.Project{}

	request := bson.M{
		"scm_uri":  scmURI,
		"scm_type": scmType,
	}
	err := c.Database.C("projects").Find(request).One(&result)
	fmt.Printf("retrieve project: %#v", result)
	return result, err
}

func (c *mongoConnector) GetProjectById(id string) (bazooka.Project, error) {
	result := bazooka.Project{}

	request := bson.M{
		"id": id,
	}
	err := c.Database.C("projects").Find(request).One(&result)
	fmt.Printf("retrieve project: %#v", result)
	return result, err
}

func (c *mongoConnector) GetProjects() ([]bazooka.Project, error) {
	result := []bazooka.Project{}

	err := c.Database.C("projects").Find(bson.M{}).All(&result)
	fmt.Printf("retrieve projects: %#v", result)
	return result, err
}

func (c *mongoConnector) AddProject(project *bazooka.Project) error {
	i := bson.NewObjectId()
	project.ID = i.Hex()

	fmt.Printf("add project: %#v", project)
	err := c.Database.C("projects").Insert(project)

	return err
}

func (c *mongoConnector) AddJob(job *bazooka.Job) error {
	fmt.Printf("add job: %#v", job)
	err := c.Database.C("jobs").Insert(job)

	return err
}

func (c *mongoConnector) GetJobByID(id string) (bazooka.Job, error) {
	result := bazooka.Job{}

	request := bson.M{
		"id": id,
	}
	err := c.Database.C("jobs").Find(request).One(&result)
	fmt.Printf("retrieve job: %#v", result)
	return result, err
}

func (c *mongoConnector) GetJobs() ([]bazooka.Job, error) {
	result := []bazooka.Job{}

	err := c.Database.C("jobs").Find(bson.M{}).All(&result)
	fmt.Printf("retrieve jobs: %#v", result)
	return result, err
}

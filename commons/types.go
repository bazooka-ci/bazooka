package bazooka

import (
	"time"
)

type Project struct {
	ScmType string `bson:"scm_type" json:"scm_type"`
	ScmURI  string `bson:"scm_uri" json:"scm_uri"`
	Name    string `bson:"name" json:"name"`
	ID      string `bson:"id" json:"id"`
}

type StartJob struct {
	ScmReference string `json:"reference"`
}

type JobStatus uint8

const (
	JOB_SUCCESS JobStatus = iota + 1
	JOB_FAILED
	JOB_ERRORED
)

type Job struct {
	ID              string    `bson:"id" json:"id"`
	ProjectID       string    `bson:"project_id" json:"project_id"`
	OrchestrationID string    `bson:"orchestration_id" json:"orchestration_id"`
	Started         time.Time `bson:"started" json:"started"`
	Completed       time.Time `bson:"completed" json:"completed"`
	Status          JobStatus `bson:"status" json:"status"`
}

type ScmFetcher struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	ImageName   string `bson:"image_name" json:"image_name"`
	ID          string `bson:"id" json:"id"`
}

type Config struct {
	Language      string   `yaml:"language"`
	BeforeInstall []string `yaml:"before_install,omitempty"`
	Install       []string `yaml:"install,omitempty"`
	BeforeScript  []string `yaml:"before_script,omitempty"`
	Script        []string `yaml:"script,omitempty"`
	AfterScript   []string `yaml:"after_script,omitempty"`
	AfterSuccess  []string `yaml:"after_success,omitempty"`
	AfterFailure  []string `yaml:"after_failure,omitempty"`
	Services      []string `yaml:"services,omitempty"`
	Env           []string `yaml:"env,omitempty"`
	FromImage     string   `yaml:"from"`
}

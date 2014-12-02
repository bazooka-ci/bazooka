package bazooka

import (
	"fmt"
	"time"
)

type Project struct {
	ScmType string `bson:"scm_type" json:"scm_type"`
	ScmURI  string `bson:"scm_uri" json:"scm_uri"`
	Name    string `bson:"name" json:"name"`
	ID      string `bson:"id" json:"id"`
}

type Variant struct {
	Status     JobStatus `bson:"status" json:"status"`
	Started    time.Time `bson:"started" json:"started"`
	Completed  time.Time `bson:"completed" json:"completed"`
	BuildImage string    `bson:"image" json:"image"`
	JobID      string    `bson:"job_id" json:"job_id"`
	Number     int       `bson:"number" json:"number"`
	ID         string    `bson:"id" json:"id"`
}

type StartJob struct {
	ScmReference string `json:"reference"`
}

type JobStatus string

const (
	JOB_SUCCESS JobStatus = "SUCCESS"
	JOB_FAILED            = "FAILED"
	JOB_ERRORED           = "ERRORED"
	JOB_RUNNING           = "RUNNING"
)

type Job struct {
	ID              string      `bson:"id" json:"id"`
	ProjectID       string      `bson:"project_id" json:"project_id"`
	OrchestrationID string      `bson:"orchestration_id" json:"orchestration_id"`
	Started         time.Time   `bson:"started" json:"started"`
	Completed       time.Time   `bson:"completed" json:"completed"`
	Status          JobStatus   `bson:"status" json:"status"`
	SCMMetadata     SCMMetadata `bson:"scm_metadata" json:"scm_metadata"`
}

type LogEntry struct {
	ID        string    `bson:"id" json:"id"`
	Message   string    `bson:"msg" json:"msg"`
	Time      time.Time `bson:"time" json:"time"`
	ProjectID string    `bson:"project_id" json:"project_id"`
	JobID     string    `bson:"job_id" json:"job_id"`
	VariantID string    `bson:"variant_id" json:"variant_id"`
	Image     string    `bson:"image" json:"image"`
}

type SCMMetadata struct {
	Reference string   `bson:"reference" json:"reference" yaml:"reference"`
	CommitID  string   `bson:"commit_id" json:"commit_id" yaml:"commit_id"`
	Author    Person   `bson:"author" json:"author" yaml:"author"`
	Date      YamlTime `bson:"time" json:"date" yaml:"date"`
	Message   string   `bson:"message" json:"message" yaml:"message"`
}

type Person struct {
	Name  string `bson:"name" json:"name" yaml:"name"`
	Email string `bson:"email" json:"email" yaml:"email"`
}

type ScmFetcher struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	ImageName   string `bson:"image_name" json:"image_name"`
	ID          string `bson:"id" json:"id"`
}

type Config struct {
	Language      string       `yaml:"language"`
	Setup         []string     `yaml:"setup,omitempty"`
	BeforeInstall []string     `yaml:"before_install,omitempty"`
	Install       []string     `yaml:"install,omitempty"`
	BeforeScript  []string     `yaml:"before_script,omitempty"`
	Script        []string     `yaml:"script,omitempty"`
	AfterScript   []string     `yaml:"after_script,omitempty"`
	AfterSuccess  []string     `yaml:"after_success,omitempty"`
	AfterFailure  []string     `yaml:"after_failure,omitempty"`
	Services      []string     `yaml:"services,omitempty"`
	Env           []string     `yaml:"env,omitempty"`
	FromImage     string       `yaml:"from"`
	Matrix        ConfigMatrix `yaml:"matrix,omitempty"`
}

type ConfigMatrix struct {
	Exclude []map[string]interface{} `yaml:"exclude,omitempty"`
}

type YamlTime struct {
	time.Time
	f string
}

func (t YamlTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	timeAsString := ""
	if err := unmarshal(&timeAsString); err != nil {
		return err
	}

	test, err := time.Parse("Mon Jan 02 15:04:05 2006 -0700", timeAsString)
	if err == nil {
		t = YamlTime{
			test,
			timeAsString,
		}
		fmt.Printf("Parsed time is %v\n", t)
		return nil
	}
	return err
}

func (j YamlTime) format() string {
	fmt.Printf("###########FORMAT\n")
	return j.Time.Format(j.f)
}

func (j YamlTime) MarshalText() ([]byte, error) {
	fmt.Printf("###########MT\n")
	return []byte(j.format()), nil
}

func (j YamlTime) MarshalJSON() ([]byte, error) {
	fmt.Printf("###########MJ\n")
	return []byte(`"` + j.format() + `"`), nil
}

func (j YamlTime) MarshalYaml() ([]byte, error) {
	fmt.Printf("###########MY\n")
	return []byte(j.format()), nil
}

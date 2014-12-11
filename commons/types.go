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
	Status     JobStatus     `bson:"status" json:"status"`
	Started    time.Time     `bson:"started" json:"started"`
	Completed  time.Time     `bson:"completed" json:"completed"`
	BuildImage string        `bson:"image" json:"image"`
	JobID      string        `bson:"job_id" json:"job_id"`
	Number     int           `bson:"number" json:"number"`
	ID         string        `bson:"id" json:"id"`
	Metas      *VariantMetas `bson:"metas" json:"metas"`
}

type VariantMetas []*VariantMeta

type VariantMeta struct {
	Kind        MetaKind
	Name, Value string
}

type MetaKind string

const (
	META_ENV  MetaKind = "env"
	META_LANG          = "lang"
)

func (ms *VariantMetas) Append(m *VariantMeta) {
	*ms = append(*ms, m)
}
func (ms *VariantMetas) Len() int      { return len(*ms) }
func (ms *VariantMetas) Swap(i, j int) { (*ms)[i], (*ms)[j] = (*ms)[j], (*ms)[i] }
func (ms *VariantMetas) Less(i, j int) bool {
	a, b := (*ms)[i], (*ms)[j]
	if a.Kind == b.Kind {
		return a.Name < b.Name
	}
	return a.Kind == META_LANG
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
	Origin    string   `bson:"origin" json:"origin" yaml:"origin"`
	Reference string   `bson:"reference" json:"reference" yaml:"reference"`
	CommitID  string   `bson:"commit_id" json:"commit_id" yaml:"commit_id"`
	Author    Person   `bson:"author" json:"author" yaml:"author"`
	Committer Person   `bson:"committer" json:"committer" yaml:"committer"`
	Date      YamlTime `bson:"time" json:"date" yaml:"date"`
	Message   string   `bson:"message" json:"message" yaml:"message"`
}

type Person struct {
	Name  string `bson:"name" json:"name" yaml:"name"`
	Email string `bson:"email" json:"email" yaml:"email"`
}

type Image struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Image       string `bson:"image" json:"image"`
	ID          string `bson:"id" json:"id"`
}

type User struct {
	ID       string `bson:"id" json:"id"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type Config struct {
	Language      string       `yaml:"language"`
	Setup         Commands     `yaml:"setup,omitempty"`
	BeforeInstall Commands     `yaml:"before_install,omitempty"`
	Install       Commands     `yaml:"install,omitempty"`
	BeforeScript  Commands     `yaml:"before_script,omitempty"`
	Script        Commands     `yaml:"script,omitempty"`
	AfterScript   Commands     `yaml:"after_script,omitempty"`
	AfterSuccess  Commands     `yaml:"after_success,omitempty"`
	AfterFailure  Commands     `yaml:"after_failure,omitempty"`
	Services      []string     `yaml:"services,omitempty"`
	Env           []string     `yaml:"env,omitempty"`
	FromImage     string       `yaml:"from"`
	Matrix        ConfigMatrix `yaml:"matrix,omitempty"`
}

type Commands []string

func (c *Commands) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	switch convCmd := raw.(type) {
	case string:
		*c = []string{convCmd}
		return nil
	case []interface{}:
		*c = make([]string, len(convCmd))
		for i, rawCmd := range convCmd {
			cmd, ok := rawCmd.(string)
			if !ok {
				return fmt.Errorf("Command list (install, script, ...) can only contain strings")
			}
			(*c)[i] = cmd
		}
		return nil
	default:
		return fmt.Errorf("Commands (install, script, ...) can be either a tring or a list of strings")
	}
}

type ConfigMatrix struct {
	Exclude []map[string]interface{} `yaml:"exclude,omitempty"`
}

type YamlTime struct {
	time.Time
}

func (t *YamlTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	timeAsString := ""
	if err := unmarshal(&timeAsString); err != nil {
		return err
	}
	if len(timeAsString) == 0 {
		return nil
	}

	timeFormats := []string{
		time.ANSIC,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		"Mon Jan 2 15:04:05 2006 -0700",
	}

	for _, timeFormat := range timeFormats {
		test, err := time.Parse(timeFormat, timeAsString)
		if err == nil {
			*t = YamlTime{
				test,
			}
			return nil
		}
	}

	return fmt.Errorf("Unable to parse time %v", timeAsString)
}

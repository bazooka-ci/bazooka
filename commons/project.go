package bazooka

import "time"

type Project struct {
	ScmType    string            `bson:"scm_type" json:"scm_type" validate:"required"`
	ScmURI     string            `bson:"scm_uri" json:"scm_uri" validate:"required"`
	Name       string            `bson:"name" json:"name" validate:"required"`
	ID         string            `bson:"id" json:"id"`
	HookKey    string            `bson:"hook_key" json:"hook_key"`
	JobCounter int               `bson:"job_counter" json:"job_counter"`
	Config     map[string]string `bson:"config" json:"config"`
}

type ProjectWithStatus struct {
	*Project
	LastJob *Job `json:"last_job"`
}

type JobStatus string

const (
	JOB_SUCCESS JobStatus = "SUCCESS"
	JOB_FAILED            = "FAILED"
	JOB_ERRORED           = "ERRORED"
	JOB_RUNNING           = "RUNNING"
	JOB_PENDING           = "PENDING"
)

type Job struct {
	ID              string      `bson:"id" json:"id"`
	Number          int         `bson:"number" json:"number"`
	ProjectID       string      `bson:"project_id" json:"project_id"`
	OrchestrationID string      `bson:"orchestration_id" json:"orchestration_id"`
	Submitted       time.Time   `bson:"submitted" json:"submitted"`
	Started         time.Time   `bson:"started" json:"started"`
	Completed       time.Time   `bson:"completed" json:"completed"`
	Status          JobStatus   `bson:"status" json:"status"`
	SCMMetadata     SCMMetadata `bson:"scm_metadata" json:"scm_metadata"`
	Parameters      []string    `bson:"parameters" json:"parameters"`
}

type Variant struct {
	Status     JobStatus     `bson:"status" json:"status"`
	Started    time.Time     `bson:"started" json:"started"`
	Completed  time.Time     `bson:"completed" json:"completed"`
	BuildImage string        `bson:"image" json:"image"`
	ProjectID  string        `bson:"project_id" json:"project_id"`
	JobID      string        `bson:"job_id" json:"job_id"`
	Number     int           `bson:"number" json:"number"`
	ID         string        `bson:"id" json:"id"`
	Metas      *VariantMetas `bson:"metas" json:"metas"`
	Artifacts  []string      `bson:"artifacts" json:"artifacts"`
}

type VariantMetas []*VariantMeta

type VariantMeta struct {
	Kind  MetaKind `bson:"kind" json:"kind"`
	Name  string   `bson:"name" json:"name"`
	Value string   `bson:"value" json:"value"`
}

type MetaKind string

const (
	META_ENV  MetaKind = "env"
	META_LANG          = "lang"
)

type StartJob struct {
	ScmReference string   `json:"reference"`
	Parameters   []string `json:"parameters"`
}

type LogEntry struct {
	ID        string    `bson:"id" json:"id"`
	Message   string    `bson:"msg" json:"msg"`
	Time      time.Time `bson:"time" json:"time"`
	Level     string    `bson:"level" json:"level"`
	Phase     string    `bson:"phase" json:"phase"`
	Command   string    `bson:"command" json:"command"`
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

type FinishData struct {
	Status    JobStatus `json:"status"`
	Time      time.Time `json:"time,omitempty"`
	Artifacts []string  `json:"artifacts,omitempty"`
}

func (ms *VariantMetas) Append(m *VariantMeta) { *ms = append(*ms, m) }
func (ms *VariantMetas) Len() int              { return len(*ms) }
func (ms *VariantMetas) Swap(i, j int)         { (*ms)[i], (*ms)[j] = (*ms)[j], (*ms)[i] }
func (ms *VariantMetas) Less(i, j int) bool {
	a, b := (*ms)[i], (*ms)[j]
	if a.Kind == b.Kind {
		return a.Name < b.Name
	}
	return a.Kind == META_LANG
}

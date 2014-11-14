package bazooka

type Project struct {
	ScmType string `bson:"scm_type" json:"scm_type"`
	ScmURI  string `bson:"scm_uri" json:"scm_uri"`
	Name    string `bson:"name" json:"name"`
	ID      string `bson:"id" json:"id"`
}

type StartJob struct {
	ScmReference string `json:"reference"`
}

type Job struct {
	ID              string `bson:"id" json:"id"`
	ProjectID       string `bson:"project_id" json:"project_id"`
	OrchestrationID string `bson:"orchestration_id" json:"orchestration_id"`
}

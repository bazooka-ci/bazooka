package bazooka

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	PrivateKey []byte
)

type Project struct {
	ScmType    string            `bson:"scm_type" json:"scm_type"`
	ScmURI     string            `bson:"scm_uri" json:"scm_uri"`
	Name       string            `bson:"name" json:"name"`
	ID         string            `bson:"id" json:"id"`
	JobCounter int               `bson:"job_counter" json:"job_counter"`
	Config     map[string]string `bson:"config" json:"config"`
}

type ProjectWithStatus struct {
	*Project
	LastJob *Job `json:"last_job"`
}

type Variant struct {
	Status     JobStatus     `bson:"status" json:"status"`
	Started    time.Time     `bson:"started" json:"started"`
	Completed  time.Time     `bson:"completed" json:"completed"`
	BuildImage string        `bson:"image" json:"image"`
	ProjectID  string        `bson:"project_id" json:"job_id"`
	JobID      string        `bson:"job_id" json:"job_id"`
	Number     int           `bson:"number" json:"number"`
	ID         string        `bson:"id" json:"id"`
	Metas      *VariantMetas `bson:"metas" json:"metas"`
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
	ScmReference string   `json:"reference"`
	Parameters   []string `json:"parameters"`
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
	Number          int         `bson:"number" json:"number"`
	ProjectID       string      `bson:"project_id" json:"project_id"`
	OrchestrationID string      `bson:"orchestration_id" json:"orchestration_id"`
	Started         time.Time   `bson:"started" json:"started"`
	Completed       time.Time   `bson:"completed" json:"completed"`
	Status          JobStatus   `bson:"status" json:"status"`
	SCMMetadata     SCMMetadata `bson:"scm_metadata" json:"scm_metadata"`
	Parameters      []string    `bson:"parameters" json:"parameters"`
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
	Language       string       `yaml:"language,omitempty"`
	Image          Images       `yaml:"image,omitempty"`
	Setup          Commands     `yaml:"setup,omitempty"`
	BeforeInstall  Commands     `yaml:"before_install,omitempty"`
	Install        Commands     `yaml:"install,omitempty"`
	BeforeScript   Commands     `yaml:"before_script,omitempty"`
	Script         Commands     `yaml:"script,omitempty"`
	AfterScript    Commands     `yaml:"after_script,omitempty"`
	AfterSuccess   Commands     `yaml:"after_success,omitempty"`
	AfterFailure   Commands     `yaml:"after_failure,omitempty"`
	Services       []string     `yaml:"services,omitempty"`
	Env            []BzkString  `yaml:"env,omitempty"`
	FromImage      string       `yaml:"from"`
	Matrix         ConfigMatrix `yaml:"matrix,omitempty"`
	Archive        Globs        `yaml:"archive,omitempty"`
	ArchiveSuccess Globs        `yaml:"archive_success,omitempty"`
	ArchiveFailure Globs        `yaml:"archive_failure,omitempty"`
}

type Images []string

func (im *Images) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	switch convCmd := raw.(type) {
	case string:
		*im = []string{convCmd}
		return nil
	case []interface{}:
		*im = make([]string, len(convCmd))
		for i, rawCmd := range convCmd {
			cmd, ok := rawCmd.(string)
			if !ok {
				return fmt.Errorf("Image list can only contain strings")
			}
			(*im)[i] = cmd
		}
		return nil
	default:
		return fmt.Errorf("Image can be either a string or a list of strings")
	}
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
		return fmt.Errorf("Commands (install, script, ...) can be either a string or a list of strings")
	}
}

type Globs []string

func (g *Globs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	switch convGlob := raw.(type) {
	case string:
		*g = []string{convGlob}
		return nil
	case []interface{}:
		*g = make([]string, len(convGlob))
		for i, rawGlob := range convGlob {
			glob, ok := rawGlob.(string)
			if !ok {
				return fmt.Errorf("Globs (archive, archive_success, archive_failure) can only contain strings")
			}
			(*g)[i] = glob
		}
		return nil
	default:
		return fmt.Errorf("Globs (archive, archive_success, archive_failure) can be either a tring or a list of strings")
	}
}

type BzkString struct {
	Name    string
	Value   string
	Secured bool
}

func (c *BzkString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	bzkStr, err := extractBzkString(raw)
	if err != nil {
		return err
	}
	*c = bzkStr
	return nil
}

var (
	_ yaml.Marshaler = BzkString{}
)

func (c BzkString) MarshalYAML() (interface{}, error) {
	merged := c.Name
	if len(c.Value) > 0 {
		merged = fmt.Sprintf("%s=%s", merged, c.Value)
	}
	if !c.Secured {
		return merged, nil
	}

	encrypted, err := encryptBzkString(merged)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"secure": encrypted,
	}, nil

}

func extractBzkString(raw interface{}) (BzkString, error) {
	switch convCmd := raw.(type) {
	case string:
		n, v := SplitNameValue(convCmd)
		return BzkString{n, v, false}, nil
	case map[interface{}]interface{}:
		if len(convCmd) > 1 {
			return BzkString{}, fmt.Errorf("BzkString should either be a go string or 'secure: <string>'")
		}
		if _, ok := convCmd["secure"]; !ok {
			return BzkString{}, fmt.Errorf("BzkString should either be a go string or 'secure: <string>'")
		}

		decrypted, err := decryptBzkString(convCmd["secure"].(string))
		if err != nil {
			return BzkString{}, fmt.Errorf("Error while trying to decrypt data, reason is: %v\n", err)
		}
		n, v := SplitNameValue(string(decrypted))
		return BzkString{n, v, true}, nil
	default:
		return BzkString{}, fmt.Errorf("BzkString should either be a go string or 'secure: <string>'")
	}
}

func decryptBzkString(str string) ([]byte, error) {
	if PrivateKey == nil {
		return nil, fmt.Errorf("PrivateKey is not set")
	}

	toDecryptDataAsHex, err := hex.DecodeString(string(str))
	if err != nil {
		return nil, fmt.Errorf("Unable to decode string as hexa data, reason is: %v\n", err)
	}

	decrypted, err := Decrypt(PrivateKey, toDecryptDataAsHex)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to decrypt data, reason is: %v\n", err)
	}
	return decrypted, nil
}

func encryptBzkString(str string) (string, error) {
	if PrivateKey == nil {
		return "", fmt.Errorf("PrivateKey is not set")
	}

	encrypted, err := Encrypt(PrivateKey, []byte(str))
	if err != nil {
		return "", fmt.Errorf("Error while trying to encrypt data, reason is: %v\n", err)
	}

	return hex.EncodeToString(encrypted), nil
}

func SplitNameValue(s string) (string, string) {
	split := strings.SplitN(string(s), "=", 2)
	value := ""
	if len(split) == 2 {
		value = split[1]
	}
	return split[0], value
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

type SSHKey struct {
	ID        string `bson:"id" json:"id"`
	Content   string `bson:"content" json:"content"`
	ProjectID string `bson:"project_id" json:"project_id"`
}

type CryptoKey struct {
	ID        string `bson:"id" json:"id"`
	Content   []byte `bson:"content" json:"content"`
	ProjectID string `bson:"project_id" json:"project_id"`
}

type StringValue struct {
	Value string `bson:"value" json:"value"`
}

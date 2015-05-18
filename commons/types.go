package bazooka

import (
	"fmt"
	"time"
)

type YamlTime struct {
	time.Time
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
	Value string `bson:"value" json:"value" validate:"required"`
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

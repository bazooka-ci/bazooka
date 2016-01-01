package bazooka

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	bazookaConfigFile = ".bazooka.yml"
	travisConfigFile  = ".travis.yml"
)

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
	Services       []Service    `yaml:"services,omitempty"`
	Env            []BzkString  `yaml:"env,omitempty"`
	FromImage      string       `yaml:"from"`
	Matrix         ConfigMatrix `yaml:"matrix,omitempty"`
	Archive        Globs        `yaml:"archive,omitempty"`
	ArchiveSuccess Globs        `yaml:"archive_success,omitempty"`
	ArchiveFailure Globs        `yaml:"archive_failure,omitempty"`
}

// Service is the representation of a a linked Docker container for the build
type Service struct {
	Image string `yaml:"image"`
	Alias string `yaml:"alias,omitempty"`
}

type Images []string

type Commands []string

type Globs []string

type ConfigMatrix struct {
	Exclude []map[string]interface{} `yaml:"exclude,omitempty"`
}

func ResolveConfigFile(source string) (string, error) {
	customPath := os.Getenv("BZK_FILE")
	if len(customPath) > 0 {
		bazookaPath := filepath.Join(source, customPath)
		exist, err := FileExists(bazookaPath)
		if err != nil {
			return "", err
		}
		if !exist {
			return "", fmt.Errorf("The custom config file %s does not exist", customPath)
		}
		return bazookaPath, nil
	}

	bazookaPath := filepath.Join(source, bazookaConfigFile)
	exist, err := FileExists(bazookaPath)
	if err != nil {
		return "", err
	}
	if exist {
		return bazookaPath, nil
	}

	travisPath := filepath.Join(source, travisConfigFile)
	exist, err = FileExists(travisPath)
	if err != nil {
		return "", err
	}
	if exist {
		return travisPath, nil
	}
	return "", errors.New("Unable to find either .bazooka.yml or .travis.yml at the root of the project")
}

func (im *Images) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	*im, err = unmarshalOneOrMany(unmarshal, "Image")
	return err
}

func (c *Commands) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	*c, err = unmarshalOneOrMany(unmarshal, "Command list (install, script, ...)")
	return err
}

func (g *Globs) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	*g, err = unmarshalOneOrMany(unmarshal, "Globs (archive, archive_success, archive_failure)")
	return err
}

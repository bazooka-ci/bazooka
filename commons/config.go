package bazooka

import (
	"errors"
	"fmt"
	"strings"
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
	Services       []string     `yaml:"services,omitempty"`
	Env            []BzkString  `yaml:"env,omitempty"`
	FromImage      string       `yaml:"from"`
	Matrix         ConfigMatrix `yaml:"matrix,omitempty"`
	Archive        Globs        `yaml:"archive,omitempty"`
	ArchiveSuccess Globs        `yaml:"archive_success,omitempty"`
	ArchiveFailure Globs        `yaml:"archive_failure,omitempty"`
}

type Images []string

type Commands []string

type Globs []string

type ConfigMatrix struct {
	Exclude []map[string]interface{} `yaml:"exclude,omitempty"`
}

func FlattenStringsEnvMap(mapp map[string]string) []string {
	res := []string{}
	for key, value := range mapp {
		res = append(res, fmt.Sprintf("%s=%s", key, value))
	}
	return res
}

func ResolveConfigFile(source string) (string, error) {
	bazookaPath := fmt.Sprintf("%s/%s", source, bazookaConfigFile)
	exist, err := FileExists(bazookaPath)
	if err != nil {
		return "", err
	}
	if exist {
		return bazookaPath, nil
	}

	travisPath := fmt.Sprintf("%s/%s", source, travisConfigFile)
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

func GetStringsEnvMap(envArray []string) map[string][]string {
	envKeyMap := make(map[string][]string)
	for _, env := range envArray {
		envSplit := strings.Split(env, "=")
		value := ""
		if len(envSplit) == 2 {
			value = envSplit[1]
		}
		envKeyMap[envSplit[0]] = append(envKeyMap[envSplit[0]], value)
	}
	return envKeyMap
}

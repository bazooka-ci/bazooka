package bazooka

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	bazookaConfigFile = ".bazooka.yml"
	travisConfigFile  = ".travis.yml"
)

func FlattenEnvMap(mapp map[BzkString]BzkString) []BzkString {
	res := []BzkString{}
	for key, value := range mapp {
		res = append(res, BzkString(fmt.Sprintf("%s=%s", key, value)))
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

func Parse(file string, object interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// TODO Add validation
	return yaml.Unmarshal(b, object)
}

func Flush(object interface{}, outputFile string) error {
	d, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFile, d, 0644)
}

func GetEnvMap(envArray []BzkString) map[string][]string {
	envKeyMap := make(map[string][]string)
	for _, env := range envArray {
		envSplit := strings.Split(string(env), "=")
		envKeyMap[envSplit[0]] = append(envKeyMap[envSplit[0]], envSplit[1])
	}
	return envKeyMap
}

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	CONFIGFILE = ".bzkcfg"
)

// TODO handle multiple server
type Config struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Auth     string `yaml:"auth"`

	Registry string `yaml:"registry"`
	Tag      string `yaml:"tag"`

	DbURL    string `yaml:"db_url"`
	QueueURL string `yaml:"queue_url"`
}

func saveConfig(authConfig *Config) error {
	confFile := path.Join(os.Getenv("HOME"), CONFIGFILE)

	authCopy := authConfig
	authCopy.Auth = encodeAuth(authCopy)
	authCopy.Username = ""
	authCopy.Password = ""

	b, err := yaml.Marshal(authCopy)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(confFile, b, 0600)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig() (*Config, error) {
	authConfig := Config{}
	confFile := path.Join(os.Getenv("HOME"), CONFIGFILE)
	if _, err := os.Stat(confFile); err != nil {
		return &authConfig, nil //missing file is not an error
	}
	b, err := ioutil.ReadFile(confFile)
	if err != nil {
		return &authConfig, err
	}

	if err := yaml.Unmarshal(b, &authConfig); err == nil {
		authConfig.Username, authConfig.Password, err = decodeAuth(authConfig.Auth)
		if err != nil {
			return &authConfig, err
		}

		return &authConfig, nil
	}
	return &authConfig, err
}

func encodeAuth(authConfig *Config) string {
	authStr := authConfig.Username + ":" + authConfig.Password
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	return string(encoded)
}

func decodeAuth(authStr string) (string, string, error) {
	decLen := base64.StdEncoding.DecodedLen(len(authStr))
	decoded := make([]byte, decLen)
	authByte := []byte(authStr)
	n, err := base64.StdEncoding.Decode(decoded, authByte)
	if err != nil {
		return "", "", err
	}
	if n > decLen {
		return "", "", fmt.Errorf("Something went wrong decoding auth config")
	}
	arr := strings.SplitN(string(decoded), ":", 2)
	if len(arr) != 2 {
		return "", "", fmt.Errorf("Invalid auth configuration file")
	}
	password := strings.Trim(arr[1], "\x00")
	return arr[0], password, nil
}

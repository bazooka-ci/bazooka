package main

import (
	lib "github.com/bazooka-ci/bazooka/commons"
)

type ConfigGolang struct {
	Base       lib.Config `yaml:",inline"`
	GoVersions []string   `yaml:"go,omitempty"`
}

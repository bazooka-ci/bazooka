package main

import (
	lib "github.com/haklop/bazooka/commons"
)

type ConfigGolang struct {
	Base       lib.Config `yaml:",inline"`
	GoVersions []string   `yaml:"go,omitempty"`
}

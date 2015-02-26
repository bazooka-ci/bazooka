package main

import (
	lib "github.com/bazooka-ci/bazooka/commons"
)

type ConfigJava struct {
	Base        lib.Config `yaml:",inline"`
	JdkVersions []string   `yaml:"jdk,omitempty"`
}

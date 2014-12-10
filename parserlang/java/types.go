package main

import (
	lib "github.com/haklop/bazooka/commons"
)

type ConfigJava struct {
	Base        lib.Config `yaml:",inline"`
	JdkVersions []string   `yaml:"jdk,omitempty"`
}

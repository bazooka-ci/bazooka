package main

import (
	lib "github.com/haklop/bazooka/commons"
)

type ConfigGolang struct {
	Language      string           `yaml:"language"`
	Setup         []string         `yaml:"setup,omitempty"`
	BeforeInstall []string         `yaml:"before_install,omitempty"`
	Install       []string         `yaml:"install,omitempty"`
	BeforeScript  []string         `yaml:"before_script,omitempty"`
	Script        []string         `yaml:"script,omitempty"`
	AfterScript   []string         `yaml:"after_script,omitempty"`
	AfterSuccess  []string         `yaml:"after_success,omitempty"`
	AfterFailure  []string         `yaml:"after_failure,omitempty"`
	Services      []string         `yaml:"services,omitempty"`
	Env           []string         `yaml:"env,omitempty"`
	GoVersions    []string         `yaml:"go,omitempty"`
	FromImage     string           `yaml:"from"`
	Matrix        lib.ConfigMatrix `yaml:"matrix,omitempty"`
}

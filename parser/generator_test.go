package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	config := &Config{
		Language:  "golang",
		FromImage: "testbazooka",
		BeforeInstall: []string{
			"gem install bundler -v 1.6.6",
			"cd source",
		},
		Install: []string{
			"travis_retry bundle _1.6.6_ install --without debug",
			"echo \"Test install\"",
		},
		BeforeScript: []string{
			"echo \"Test before script1\"",
			"echo \"Test before script2\"",
		},
		Script: []string{
			"bundle _1.6.5_ exec rake",
			"bundle _1.6.6_ exec rake",
		},
		AfterScript: []string{
			"echo \"Test after script1\"",
			"echo \"Test after script2\"",
		},
		AfterSuccess: []string{
			"echo \"Test after success1\"",
			"echo \"Test after success2\"",
		},
		AfterFailure: []string{
			"echo \"Test after failure1\"",
			"echo \"Test after failure2\"",
		},
		Env: []string{
			"TEST1=test1a",
			"TEST2=test2b",
		},
	}

	g := &Generator{
		Config:       config,
		OutputFolder: "test/generator",
	}
	err := g.GenerateDockerfile()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	breal, err := ioutil.ReadFile("../test/generator/Dockerfileexpected")
	bexpected, err := ioutil.ReadFile("../test/generator/Dockerfile0")
	assert.Equal(t, breal, bexpected)

	breal, err = ioutil.ReadFile("../test/generator/bazooka_run_expected.sh")
	bexpected, err = ioutil.ReadFile("../test/generator/bazooka_run0.sh")
	assert.Equal(t, breal, bexpected)
}

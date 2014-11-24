package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	lib "github.com/haklop/bazooka/commons"
)

type Generator struct {
	Config       *lib.Config
	OutputFolder string
	Index        int
}

func (g *Generator) GenerateDockerfile() error {
	err := os.MkdirAll(fmt.Sprintf("%s/%d", g.OutputFolder, g.Index), 0755)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("FROM %s\n\n", g.Config.FromImage))

	buffer.WriteString("ADD . /bazooka\n\n")
	buffer.WriteString("RUN chmod +x /bazooka/bazooka_run.sh\n\n")

	type buildPhase struct {
		name      string
		commands  []string
		beforeCmd []string
		runCmd    []string
	}

	phases := []*buildPhase{
		&buildPhase{
			name:      "before_install",
			commands:  g.Config.BeforeInstall,
			beforeCmd: []string{"set -ev"},
			runCmd: []string{
				"./bazooka_before_install.sh",
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "install",
			commands:  g.Config.Install,
			beforeCmd: []string{"set -ev"},
			runCmd: []string{
				"./bazooka_install.sh",
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "before_script",
			commands:  g.Config.BeforeScript,
			beforeCmd: []string{"set -ev"},
			runCmd: []string{
				"./bazooka_before_script.sh",
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "script",
			commands:  g.Config.Script,
			beforeCmd: []string{"set -v"},
			runCmd:    g.getScriptRunCmd(),
		},
		&buildPhase{
			name:      "after_success",
			commands:  g.Config.AfterSuccess,
			beforeCmd: []string{"set -v"},
			runCmd:    []string{},
		},
		&buildPhase{
			name:      "after_failure",
			commands:  g.Config.AfterFailure,
			beforeCmd: []string{"set -v"},
			runCmd:    []string{},
		},
		&buildPhase{
			name:      "after_script",
			commands:  g.Config.AfterScript,
			beforeCmd: []string{"set -v"},
		},
	}

	var bufferRun bytes.Buffer
	bufferRun.WriteString("#!/bin/bash\n")
	for _, phase := range phases {
		if len(phase.commands) != 0 {
			var buffer bytes.Buffer
			buffer.WriteString("#!/bin/bash\n\n")
			for _, action := range phase.beforeCmd {
				buffer.WriteString(fmt.Sprintf("%s\n", action))
			}
			for _, action := range phase.commands {
				buffer.WriteString(fmt.Sprintf("%s\n", action))
			}
			err = ioutil.WriteFile(fmt.Sprintf("%s/%d/bazooka_%s.sh", g.OutputFolder, g.Index, phase.name), buffer.Bytes(), 0644)
			if err != nil {
				return fmt.Errorf("Phase [%d/%s]: writing file failed: %v", g.Index, phase.name, err)
			}
			if phase.runCmd == nil {
				bufferRun.WriteString(fmt.Sprintf("./bazooka_%s.sh\n", phase.name))
			} else {
				for _, action := range phase.runCmd {
					bufferRun.WriteString(fmt.Sprintf("%s\n", action))
				}
			}
		}
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%d/bazooka_run.sh", g.OutputFolder, g.Index), bufferRun.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Phase [%d/run]: writing file failed: %v", g.Index, err)
	}

	for _, phase := range phases {
		if len(phase.commands) != 0 {
			buffer.WriteString(fmt.Sprintf("RUN chmod +x /bazooka/bazooka_%s.sh\n\n", phase.name))
		}
	}

	for _, env := range g.Config.Env {
		envSplit := strings.Split(env, "=")
		buffer.WriteString(fmt.Sprintf("ENV %s %s\n", envSplit[0], envSplit[1]))
	}

	buffer.WriteString("WORKDIR /bazooka\n\n")

	buffer.WriteString("CMD ./bazooka_run.sh\n")

	err = ioutil.WriteFile(fmt.Sprintf("%s/%d/Dockerfile", g.OutputFolder, g.Index), buffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Phase [%d/docker]: writing file failed: %v", g.Index, err)
	}

	if len(g.Config.Services) > 0 {
		var servicesBuffer bytes.Buffer
		for _, service := range g.Config.Services {
			servicesBuffer.WriteString(fmt.Sprintf("%s\n", service))
		}
		err = ioutil.WriteFile(fmt.Sprintf("%s/%d/services", g.OutputFolder, g.Index), servicesBuffer.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("Phase [%d/services]: writing file failed: %v", g.Index, err)
		}
	}
	return nil
}

func (g *Generator) getScriptRunCmd() []string {
	switch {
	case len(g.Config.AfterSuccess) != 0 && len(g.Config.AfterFailure) != 0:
		return []string{
			"./bazooka_script.sh",
			"if [[ $? != 0 ]] ; then",
			"  ./bazooka_after_failure.sh",
			"else",
			"  ./bazooka_after_success.sh",
			"fi",
		}
	case len(g.Config.AfterSuccess) != 0 && len(g.Config.AfterFailure) == 0:
		return []string{
			"./bazooka_script.sh",
			"if [[ $? == 0 ]] ; then",
			"  ./bazooka_after_success.sh",
			"fi",
		}
	case len(g.Config.AfterSuccess) == 0 && len(g.Config.AfterFailure) == 0:
		return nil
	case len(g.Config.AfterSuccess) == 0 && len(g.Config.AfterFailure) != 0:
		return []string{
			"./bazooka_script.sh",
			"if [[ $? != 0 ]] ; then",
			"  ./bazooka_after_failure.sh",
			"fi",
		}
	default:
		return nil
	}
}

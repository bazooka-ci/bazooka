package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	lib "github.com/bazooka-ci/bazooka/commons"
)

type Generator struct {
	Config       *lib.Config
	OutputFolder string
	Index        string
}

func (g *Generator) GenerateDockerfile() error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", g.OutputFolder, g.Index), 0755)
	if err != nil {
		return err
	}

	var dockerBuffer bytes.Buffer

	dockerBuffer.WriteString(fmt.Sprintf("FROM %s\n\n", g.Config.FromImage))

	envMap := lib.GetEnvMap(g.Config.Env)

	dockerBuffer.WriteString(fmt.Sprintf("COPY source %s/\n\n", envMap["BZK_BUILD_DIR"][0]))

	dockerBuffer.WriteString(fmt.Sprintf("COPY work/%s/bazooka_run.sh %s/\n", g.Index, envMap["BZK_BUILD_DIR"][0]))
	dockerBuffer.WriteString(fmt.Sprintf("RUN  chmod +x %s/bazooka_run.sh\n\n", envMap["BZK_BUILD_DIR"][0]))

	type buildPhase struct {
		name      string
		commands  []string
		beforeCmd []string
		runCmd    []string
		always    bool
	}

	phases := []*buildPhase{
		&buildPhase{
			name:      "setup",
			commands:  g.Config.Setup,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				"./bazooka_setup.sh",
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "before_install",
			commands:  g.Config.BeforeInstall,
			beforeCmd: []string{"set -e"},
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
			beforeCmd: []string{"set -e"},
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
			beforeCmd: []string{"set -e"},
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
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				"./bazooka_script.sh",
				"exitCode=$?",
			},
		},
		&buildPhase{
			name:     "archive",
			commands: archiveCommands(g.Config.Archive),
		},
		&buildPhase{
			name:     "archive_success",
			commands: archiveCommands(g.Config.ArchiveSuccess),
			runCmd: []string{
				"if [[ $exitCode == 0 ]] ; then",
				"  ./bazooka_archive_success.sh",
				"fi",
			},
		},
		&buildPhase{
			name:     "archive_failure",
			commands: archiveCommands(g.Config.ArchiveFailure),
			runCmd: []string{
				"if [[ $exitCode != 0 ]] ; then",
				"  ./bazooka_archive_failure.sh",
				"fi",
			},
		},
		&buildPhase{
			name:      "after_success",
			commands:  g.Config.AfterSuccess,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				"if [[ $exitCode == 0 ]] ; then",
				"  ./bazooka_after_success.sh",
				"fi",
			},
		},
		&buildPhase{
			name:      "after_failure",
			commands:  g.Config.AfterFailure,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				"if [[ $exitCode != 0 ]] ; then",
				"  ./bazooka_after_failure.sh",
				"fi"},
		},
		&buildPhase{
			name:      "after_script",
			commands:  g.Config.AfterScript,
			beforeCmd: []string{"set -e"},
			runCmd:    []string{},
		},
		&buildPhase{
			name:   "end",
			always: true,
			runCmd: []string{
				"exit $exitCode",
			},
		},
	}

	var bufferRun bytes.Buffer
	bufferRun.WriteString("#!/bin/bash\n")

	for _, phase := range phases {
		if len(phase.commands) > 0 {
			var phaseBuffer bytes.Buffer
			phaseBuffer.WriteString("#!/bin/bash\n\n")
			phaseBuffer.WriteString(fmt.Sprintf("echo %s\n", strconv.Quote(fmt.Sprintf("<PHASE:%s>", phase.name))))
			for _, action := range phase.beforeCmd {
				phaseBuffer.WriteString(fmt.Sprintf("%s\n", action))
			}
			for _, action := range phase.commands {
				phaseBuffer.WriteString(fmt.Sprintf("echo %s\n", strconv.Quote(fmt.Sprintf("<CMD:%s>", action))))
				phaseBuffer.WriteString(fmt.Sprintf("%s\n", action))
			}
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s/bazooka_%s.sh", g.OutputFolder, g.Index, phase.name), phaseBuffer.Bytes(), 0644)
			if err != nil {
				return fmt.Errorf("Phase [%d/%s]: writing file failed: %v", g.Index, phase.name, err)
			}

			dockerBuffer.WriteString(fmt.Sprintf("COPY work/%s/bazooka_%s.sh %s/\n", g.Index, phase.name, envMap["BZK_BUILD_DIR"][0]))
			dockerBuffer.WriteString(fmt.Sprintf("RUN  chmod +x %s/bazooka_%s.sh\n\n", envMap["BZK_BUILD_DIR"][0], phase.name))

			if len(phase.runCmd) == 0 {
				bufferRun.WriteString(fmt.Sprintf("./bazooka_%s.sh\n", phase.name))
			} else {
				for _, action := range phase.runCmd {
					bufferRun.WriteString(fmt.Sprintf("%s\n", action))
				}
			}
		} else if phase.always {
			for _, action := range phase.runCmd {
				bufferRun.WriteString(fmt.Sprintf("%s\n", action))
			}
		}
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/bazooka_run.sh", g.OutputFolder, g.Index), bufferRun.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Phase [%d/run]: writing file failed: %v", g.Index, err)
	}

	for _, env := range g.Config.Env {
		envSplit := strings.SplitN(string(env), "=", 2)
		dockerBuffer.WriteString(fmt.Sprintf("ENV  %s %s\n", envSplit[0], envSplit[1]))
	}

	dockerBuffer.WriteString(fmt.Sprintf("WORKDIR %s\n\n", envMap["BZK_BUILD_DIR"][0]))

	dockerBuffer.WriteString("CMD  ./bazooka_run.sh\n")

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/Dockerfile", g.OutputFolder, g.Index), dockerBuffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Phase [%d/docker]: writing file failed: %v", g.Index, err)
	}

	if len(g.Config.Services) > 0 {
		var servicesBuffer bytes.Buffer
		for _, service := range g.Config.Services {
			servicesBuffer.WriteString(fmt.Sprintf("%s\n", service))
		}
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s/services", g.OutputFolder, g.Index), servicesBuffer.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("Phase [%d/services]: writing file failed: %v", g.Index, err)
		}
	}
	return nil
}

func archiveCommands(globs lib.Globs) []string {
	res := make([]string, len(globs))
	for i, pat := range globs {
		res[i] = fmt.Sprintf("cp -R %s /artifacts/", pat)
	}
	return res
}

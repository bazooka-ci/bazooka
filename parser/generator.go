package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	lib "github.com/bazooka-ci/bazooka/commons"
)

type Generator struct {
	config       *lib.Config
	outputFolder string
	index        string
}

func (g *Generator) GenerateDockerfile() error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", g.outputFolder, g.index), 0755)
	if err != nil {
		return err
	}

	var dockerBuffer bytes.Buffer

	dockerBuffer.WriteString(fmt.Sprintf("FROM %s\n\n", g.config.FromImage))

	envMap := lib.GetEnvMap(g.config.Env)

	bzkBuildDir := envMap["BZK_BUILD_DIR"][0].Value

	dockerBuffer.WriteString(fmt.Sprintf("COPY source %s/\n\n", bzkBuildDir))

	dockerBuffer.WriteString(fmt.Sprintf("COPY work/%s/bazooka_run.sh %s/\n", g.index, bzkBuildDir))
	dockerBuffer.WriteString(fmt.Sprintf("RUN  chmod +x %s/bazooka_run.sh\n\n", bzkBuildDir))

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
			commands:  g.config.Setup,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				fmt.Sprintf("%s/bazooka_setup.sh", bzkBuildDir),
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "before_install",
			commands:  g.config.BeforeInstall,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				fmt.Sprintf("%s/bazooka_before_install.sh", bzkBuildDir),
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "install",
			commands:  g.config.Install,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				fmt.Sprintf("%s/bazooka_install.sh", bzkBuildDir),
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "before_script",
			commands:  g.config.BeforeScript,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				fmt.Sprintf("%s/bazooka_before_script.sh", bzkBuildDir),
				"rc=$?",
				"if [[ $rc != 0 ]] ; then",
				"    exit 42",
				"fi",
			},
		},
		&buildPhase{
			name:      "script",
			commands:  g.config.Script,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				fmt.Sprintf("%s/bazooka_script.sh", bzkBuildDir),
				"exitCode=$?",
			},
		},
		&buildPhase{
			name:     "archive",
			commands: archiveCommands(g.config.Archive),
		},
		&buildPhase{
			name:     "archive_success",
			commands: archiveCommands(g.config.ArchiveSuccess),
			runCmd: []string{
				"if [[ $exitCode == 0 ]] ; then",
				fmt.Sprintf("  %s/bazooka_archive_success.sh", bzkBuildDir),
				"fi",
			},
		},
		&buildPhase{
			name:     "archive_failure",
			commands: archiveCommands(g.config.ArchiveFailure),
			runCmd: []string{
				"if [[ $exitCode != 0 ]] ; then",
				fmt.Sprintf("  %s/bazooka_archive_failure.sh", bzkBuildDir),
				"fi",
			},
		},
		&buildPhase{
			name:      "after_success",
			commands:  g.config.AfterSuccess,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				"if [[ $exitCode == 0 ]] ; then",
				fmt.Sprintf("  %s/bazooka_after_success.sh", bzkBuildDir),
				"fi",
			},
		},
		&buildPhase{
			name:      "after_failure",
			commands:  g.config.AfterFailure,
			beforeCmd: []string{"set -e"},
			runCmd: []string{
				"if [[ $exitCode != 0 ]] ; then",
				fmt.Sprintf("  %s/bazooka_after_failure.sh", bzkBuildDir),
				"fi"},
		},
		&buildPhase{
			name:      "after_script",
			commands:  g.config.AfterScript,
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
			bufferRun.WriteString("source /.bzkenv\n")
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
			phaseBuffer.WriteString(`echo "cd \"$(pwd)\"" > /.bzkenv`)
			phaseBuffer.WriteString("\n")
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s/bazooka_%s.sh", g.outputFolder, g.index, phase.name), phaseBuffer.Bytes(), 0644)
			if err != nil {
				return fmt.Errorf("Phase [%s/%s]: writing file failed: %v", g.index, phase.name, err)
			}

			dockerBuffer.WriteString(fmt.Sprintf("COPY work/%s/bazooka_%s.sh %s/\n", g.index, phase.name, bzkBuildDir))
			dockerBuffer.WriteString(fmt.Sprintf("RUN  chmod +x %s/bazooka_%s.sh\n\n", bzkBuildDir, phase.name))

			if len(phase.runCmd) == 0 {
				bufferRun.WriteString(fmt.Sprintf("%s/bazooka_%s.sh\n", bzkBuildDir, phase.name))
			} else {
				for _, action := range phase.runCmd {
					bufferRun.WriteString(fmt.Sprintf("%s\n", action))
				}
			}
		} else if phase.always {
			bufferRun.WriteString("source /.bzkenv\n")
			for _, action := range phase.runCmd {
				bufferRun.WriteString(fmt.Sprintf("%s\n", action))
			}
		}
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/bazooka_run.sh", g.outputFolder, g.index), bufferRun.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Phase [%s/run]: writing file failed: %v", g.index, err)
	}

	for _, env := range g.config.Env {
		dockerBuffer.WriteString(fmt.Sprintf("ENV  %s %s\n", env.Name, env.Value))
	}

	dockerBuffer.WriteString(fmt.Sprintf("WORKDIR %s\n\n", bzkBuildDir))

	dockerBuffer.WriteString(fmt.Sprintf("CMD  %s/bazooka_run.sh\n", bzkBuildDir))

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/Dockerfile", g.outputFolder, g.index), dockerBuffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Phase [%s/docker]: writing file failed: %v", g.index, err)
	}

	if len(g.config.Services) > 0 {
		err = lib.Flush(g.config.Services, fmt.Sprintf("%s/%s/services", g.outputFolder, g.index))
		if err != nil {
			return fmt.Errorf("Phase [%s/services]: writing file failed: %v", g.index, err)
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

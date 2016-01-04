package main

import (
	"fmt"
	"os"
	"text/template"

	lib "github.com/bazooka-ci/bazooka/commons"
)

type Generator struct {
	Config       *lib.Config
	OutputFolder string
	Index        string
}

type TemplateValues struct {
	Generator   *Generator
	BzkBuildDir string
	Phases      []*BuildPhase
}

type BuildPhase struct {
	Name               string
	Commands           []string
	ContinueOnCmdError bool
}

func (g *Generator) GenerateDockerfile() error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", g.OutputFolder, g.Index), 0755)
	if err != nil {
		return err
	}

	phases := []*BuildPhase{
		&BuildPhase{
			Name:     "before_install",
			Commands: g.Config.BeforeInstall,
		},
		&BuildPhase{
			Name:     "install",
			Commands: g.Config.Install,
		},
		&BuildPhase{
			Name:     "before_script",
			Commands: g.Config.BeforeScript,
		},
		&BuildPhase{
			Name:     "script",
			Commands: g.Config.Script,
		},
		&BuildPhase{
			Name:               "archive",
			Commands:           archiveCommands(g.Config.Archive),
			ContinueOnCmdError: true,
		},
		&BuildPhase{
			Name:               "archive_success",
			Commands:           archiveCommands(g.Config.ArchiveSuccess),
			ContinueOnCmdError: true,
		},
		&BuildPhase{
			Name:               "archive_failure",
			Commands:           archiveCommands(g.Config.ArchiveFailure),
			ContinueOnCmdError: true,
		},
		&BuildPhase{
			Name:     "after_success",
			Commands: g.Config.AfterSuccess,
		},
		&BuildPhase{
			Name:     "after_failure",
			Commands: g.Config.AfterFailure,
		},
		&BuildPhase{
			Name:     "after_script",
			Commands: g.Config.AfterScript,
		},
	}

	templateValues := &TemplateValues{
		Generator:   g,
		BzkBuildDir: lib.GetEnvMap(g.Config.Env)["BZK_BUILD_DIR"][0].Value,
		Phases:      phases,
	}

	err = writeTemplate(templateValues, "/template/Dockerfile", fmt.Sprintf("%s/%s/Dockerfile", g.OutputFolder, g.Index))
	if err != nil {
		return err
	}

	err = writeTemplate(templateValues, "/template/bazooka_run.sh", fmt.Sprintf("%s/%s/bazooka_run.sh", g.OutputFolder, g.Index))
	if err != nil {
		return err
	}

	for _, phase := range phases {
		err = writeTemplate(phase, "/template/bazooka_phase.sh", fmt.Sprintf("%s/%s/bazooka_%s.sh", g.OutputFolder, g.Index, phase.Name))
		if err != nil {
			return err
		}
	}

	if len(g.Config.Services) > 0 {
		err = lib.Flush(g.Config.Services, fmt.Sprintf("%s/%s/services", g.OutputFolder, g.Index))
		if err != nil {
			return fmt.Errorf("Phase [%s/services]: writing file failed: %v", g.Index, err)
		}
	}
	return nil
}

func archiveCommands(globs lib.Globs) []string {
	res := make([]string, len(globs))
	for i, pat := range globs {
		res[i] = fmt.Sprintf("cp -R \"%s\" /artifacts/", pat)
	}
	return res
}

func writeTemplate(t interface{}, templateFile, outputFile string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("Error parsing template file %s: %v", templateFile, err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("Error creating file %s: %v", outputFile, err)
	}

	err = tmpl.Execute(file, t)
	if err != nil {
		return fmt.Errorf("Error executing template file %s: %v", templateFile, err)
	}

	return nil
}

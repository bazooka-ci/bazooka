package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:        "project",
			Usage:       "Actions on projects",
			Subcommands: []cli.Command{createProjectCommand(), listProjectsCommand()},
		}, {
			Name:        "job",
			Usage:       "Actions on projects",
			Subcommands: []cli.Command{startJobCommand(), listJobsCommand(), jobLogCommand()},
		}, {
			Name:        "variant",
			Usage:       "Actions on variants",
			Subcommands: []cli.Command{listVariantsCommand(), variantLogCommand()},
		}, {
			Name:        "image",
			Usage:       "Actions on images",
			Subcommands: []cli.Command{listImagesCommand(), setImageCommand()},
		},
	}
	app.Run(os.Args)
}

const (
	idExcerptLen = 8
)

func idExcerpt(id string) string {
	if len(id) > idExcerptLen {
		return id[:idExcerptLen]
	}
	return id
}

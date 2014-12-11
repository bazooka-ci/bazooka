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
		}, {
			Name:        "user",
			Usage:       "Actions on users",
			Subcommands: []cli.Command{listUsersCommand(), createUserCommand()},
		}, {
			Name:   "login",
			Usage:  "Log in to a Bazooka server",
			Action: login,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "bazooka-uri",
					Value:  "http://localhost:3000",
					Usage:  "URI for the bazooka server",
					EnvVar: "BZK_URI",
				},
				cli.StringFlag{
					Name:   "email",
					Usage:  "User email",
					EnvVar: "BZK_USER_EMAIL",
				},
				cli.StringFlag{
					Name:   "password",
					Usage:  "User password",
					EnvVar: "BZK_USER_PASSWORD",
				},
			},
		}, {
			Name:   "run",
			Usage:  "Run bazooka",
			Action: run,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "upgrade",
					Usage: "Upgrade Bazooka to the latest version by pulling new images from the registry",
				},
				cli.BoolFlag{
					Name:  "update",
					Usage: "Update Bazooka by destroying and recreating initial containers",
				},
				cli.BoolFlag{
					Name:  "restart",
					Usage: "Restart Bazooka if already running",
				},
				cli.StringFlag{
					Name:   "registry",
					Usage:  "Custom registry to get Bazooka Docker images",
					EnvVar: "BZK_REGISTRY",
				},
				cli.StringFlag{
					Name:   "home",
					Usage:  "Home Folder for Bazooka to work",
					EnvVar: "BZK_HOME",
				},
				cli.StringFlag{
					Name:   "docker-sock",
					Usage:  "Location of the Docker unix socket, usually /var/run/docker.sock",
					Value:  "/var/run/docker.sock",
					EnvVar: "BZK_DOCKERSOCK",
				},
				cli.StringFlag{
					Name:   "scm-key",
					Usage:  "Location of the private SSH Key Bazooka will use for SCM Fetch",
					EnvVar: "BZK_SCM_KEYFILE",
				},
			},
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

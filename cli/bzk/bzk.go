package main

import (
	"os"

	"github.com/jawher/mow.cli"
)

var (
	app       = cli.App("bzk", "Bazooka CI client")
	bzkApiUrl = app.String(cli.StringOpt{
		Name:   "u bazooka-url",
		Desc:   "URL for the bazooka server",
		EnvVar: "BZK_API_URL",
	})
)

func main() {
	app.Command("project", "Actions on projects", func(cmd *cli.Cmd) {
		cmd.Command("list", "List bazooka projects", listProjectsCommand)
		cmd.Command("create", "Create a new bazooka project", createProjectCommand)
		cmd.Command("show", "Display specific information about a bazooka project", showProjectCommand)
		cmd.Command("config", "View or modify a bazooka project configuration", func(cfgCmd *cli.Cmd) {
			cfgCmd.Command("list", "List full project configuration", listProjectConfigCommand)
			cfgCmd.Command("get", "Get a specific project configuration key", getProjectConfigKeyCommand)
			cfgCmd.Command("set", "Set a specific project configuration key", setProjectConfigKeyCommand)
			cfgCmd.Command("unset", "Delete a specific project configuration key", unsetProjectConfigKeyCommand)
		})
	})

	app.Command("job", "Actions on jobs", func(cmd *cli.Cmd) {
		cmd.Command("list", "List jobs associated with a project", listJobsCommand)
		cmd.Command("start", "Start a new bazooka job on a project", startJobCommand)
		cmd.Command("log", "View a job log", jobLogCommand)
	})

	app.Command("variant", "Actions on job variants", func(cmd *cli.Cmd) {
		cmd.Command("list", "List variants associated with a job", listVariantsCommand)
		cmd.Command("log", "View a variant log", variantLogCommand)
	})

	app.Command("key", "Actions on projects keys", func(cmd *cli.Cmd) {
		cmd.Command("get", "list Keys for the bazooka project", getKeyCommand)
		cmd.Command("set", "Add SSH Key for the bazooka project", setKeyCommand)
	})

	app.Command("image", "Actions on images", func(cmd *cli.Cmd) {
		cmd.Command("list", "List the registered docker images", listImagesCommand)
		cmd.Command("get", "View a registered docker image", getImageCommand)
		cmd.Command("register", "Register a docker image", setImageCommand)
	})

	app.Command("user", "Actions on users", func(cmd *cli.Cmd) {
		cmd.Command("list", "List bazooka users", listUsersCommand)
		cmd.Command("create", "Create a new bazooka user", createUserCommand)
	})

	app.Command("encrypt", "Encrypt some data", encryptData)

	app.Command("service", "Manage bazooka service (start, stop, status, upgrade...)", func(cmd *cli.Cmd) {
		cmd.Command("start", "Start bazooka", startService)
		cmd.Command("restart", "Restart bazooka", restartService)
		cmd.Command("upgrade", "Upgrade bazooka to the latest version", upgradeService)
		cmd.Command("stop", "Stop bazooka", stopService)
		cmd.Command("status", "Get bazooka status", statusService)
	})

	app.Command("login", "Log in to the bazooka server", login)

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

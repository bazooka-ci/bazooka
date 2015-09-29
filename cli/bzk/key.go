package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jawher/mow.cli"
)

func setKeyCommand(cmd *cli.Cmd) {
	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "Project id",
	})
	scmKey := cmd.String(cli.StringArg{
		Name: "SCM_KEY_PATH",
		Desc: "The absolute path to the SCM key",
	})

	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		if err := client.Project.Key.Set(*pid, *scmKey); err != nil {
			log.Fatal(err)
		}
	}
}

func getKeyCommand(cmd *cli.Cmd) {
	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "Project id",
	})

	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.Project.Key.Get(*pid)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "PROJECT ID\n")
		fmt.Fprintf(w, "%s\n", idExcerpt(res.ProjectID))
		w.Flush()
	}
}

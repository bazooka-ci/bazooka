package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jawher/mow.cli"
)

func addKeyCommand(cmd *cli.Cmd) {
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
		res, err := client.Project.AddKey(*pid, *scmKey)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "PROJECT ID\n")
		fmt.Fprintf(w, "%s\n", idExcerpt(res.ProjectID))
		w.Flush()
	}
}

func updateKeyCommand(cmd *cli.Cmd) {
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
		res, err := client.Project.UpdateKey(*pid, *scmKey)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "PROJECT ID\n")
		fmt.Fprintf(w, "%s\n", idExcerpt(res.ProjectID))
		w.Flush()
	}
}

func listKeysCommand(cmd *cli.Cmd) {
	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "Project id",
	})

	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.Project.Keys(*pid)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "PROJECT ID\n")
		for _, item := range res {
			fmt.Fprintf(w, "%s\n", idExcerpt(item.ProjectID))
		}
		w.Flush()
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jawher/mow.cli"
)

func createProjectCommand(cmd *cli.Cmd) {
	cmd.Spec = "NAME SCM_TYPE SCM_URI [SCM_KEY]"

	name := cmd.String(cli.StringArg{
		Name: "NAME",
		Desc: "the project name",
	})
	scmType := cmd.String(cli.StringArg{
		Name: "SCM_TYPE",
		Desc: "one of the supported scm types (git, svm ...)",
	})
	scmUri := cmd.String(cli.StringArg{
		Name: "SCM_URI",
		Desc: "the project clone url",
	})
	scmKey := cmd.String(cli.StringArg{
		Name: "SCM_KEY",
		Desc: "the project SCM key",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.CreateProject(*name, *scmType, *scmUri)
		if err != nil {
			log.Fatal(err)
		}
		if len(*scmKey) > 0 {
			_, err = client.AddKey(res.ID, *scmKey)
			if err != nil {
				log.Fatal(err)
			}
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "PROJECT ID\tNAME\tSCM TYPE\tSCM URI\n")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", idExcerpt(res.ID), res.Name, res.ScmType, res.ScmURI)
		w.Flush()
	}
}

func listProjectsCommand(cmd *cli.Cmd) {
	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.ListProjects()
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "PROJECT ID\tNAME\tSCM TYPE\tSCM URI\n")
		for _, item := range res {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", idExcerpt(item.ID), item.Name, item.ScmType, item.ScmURI)
		}
		w.Flush()
	}
}

func listProjectConfigCommand(cmd *cli.Cmd) {
	cmd.Spec = "PROJECT_ID"

	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "the project id",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.GetProjectConfig(*pid)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "KEY\tVALUE\n")

		for k, v := range res {
			fmt.Fprintf(w, "%s\t%s\n", k, v)
		}
		w.Flush()
	}
}

func getProjectConfigKeyCommand(cmd *cli.Cmd) {
	cmd.Spec = "PROJECT_ID KEY"

	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "the project id",
	})

	key := cmd.String(cli.StringArg{
		Name: "KEY",
		Desc: "the configuration key",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.GetProjectConfig(*pid)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "KEY\tVALUE\n")

		v, found := res[*key]
		if !found {
			v = "<not set>"
		}

		fmt.Fprintf(w, "%s\t%s\n", *key, v)
		w.Flush()
	}
}

func setProjectConfigKeyCommand(cmd *cli.Cmd) {
	cmd.Spec = "PROJECT_ID KEY VALUE"

	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "the project id",
	})

	key := cmd.String(cli.StringArg{
		Name: "KEY",
		Desc: "the configuration key",
	})

	value := cmd.String(cli.StringArg{
		Name: "VALUE",
		Desc: "the configuration value",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		if err := client.SetProjectConfigKey(*pid, *key, *value); err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "KEY\tVALUE\n")
		fmt.Fprintf(w, "%s\t%s\n", *key, *value)
		w.Flush()
	}
}

func unsetProjectConfigKeyCommand(cmd *cli.Cmd) {
	cmd.Spec = "PROJECT_ID KEY"

	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "the project id",
	})

	key := cmd.String(cli.StringArg{
		Name: "KEY",
		Desc: "the configuration key",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		if err := client.UnsetProjectConfigKey(*pid, *key); err != nil {
			log.Fatal(err)
		}
	}
}

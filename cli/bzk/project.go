package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func createProjectCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create a new Project on Bazooka",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "bazooka-uri",
				Value:  "http://localhost:3000",
				Usage:  "URI for the bazooka server",
				EnvVar: "BZK_URI",
			},
			cli.StringFlag{
				Name:  "build",
				Value: "scm",
				Usage: "Type of build for Bazooka, eg. scm for `SCM based build`",
			},
			cli.StringFlag{
				Name:   "name",
				Usage:  "Name of the project",
				EnvVar: "BZK_SCM_NAME",
			},
			cli.StringFlag{
				Name:   "scm",
				Value:  "git",
				Usage:  "Type of SCM for the SCM based-build",
				EnvVar: "BZK_SCM",
			},
			cli.StringFlag{
				Name:   "scm-uri",
				Usage:  "URI of your SCM Project",
				EnvVar: "BZK_SCM_URI",
			},
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			res, err := client.CreateProject(c.String("name"), c.String("scm"), c.String("scm-uri"))
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "PROJECT ID\tNAME\tSCM TYPE\tSCM URI\n")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", res.ID, res.Name, res.ScmType, res.ScmURI)
			w.Flush()
		},
	}
}

func listProjectsCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List bazooka projects",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "bazooka-uri",
				Value:  "http://localhost:3000",
				Usage:  "URI for the bazooka server",
				EnvVar: "BZK_URI",
			},
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
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
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", item.ID, item.Name, item.ScmType, item.ScmURI)
			}
			w.Flush()
		},
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	create := cli.Command{
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
			fmt.Fprint(w, "JOB ID\tPROJECT ID\tORCHESTRATION ID\n")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", res.ID, res.Name, res.ScmType, res.ScmURI)
			w.Flush()
		},
	}

	list := cli.Command{
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

	startJob := cli.Command{
		Name:  "start-job",
		Usage: "Start Job for the bazooka project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "bazooka-uri",
				Value:  "http://localhost:3000",
				Usage:  "URI for the bazooka server",
				EnvVar: "BZK_URI",
			},
			cli.StringFlag{
				Name:  "project-id",
				Usage: "ID of the project to buiod",
			},
			cli.StringFlag{
				Name:  "scm-ref",
				Value: "master",
				Usage: "SCM Reference to build",
			},
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			res, err := client.StartJob(c.String("project-id"), c.String("scm-ref"))
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "JOB ID\tPROJECT ID\tORCHESTRATION ID\n")
			fmt.Fprintf(w, "%s\t%s\t%s\t\n", res.ID, res.ProjectID, res.OrchestrationID)
			w.Flush()
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "project",
			Usage:       "Actions on projects",
			Subcommands: []cli.Command{create, list, startJob},
		},
	}
	app.Run(os.Args)
}

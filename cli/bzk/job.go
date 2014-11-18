package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func startJobCommand() cli.Command {
	return cli.Command{
		Name:  "start",
		Usage: "Start Job for the bazooka project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "bazooka-uri",
				Value:  "http://localhost:3000",
				Usage:  "URI for the bazooka server",
				EnvVar: "BZK_URI",
			},
			cli.StringFlag{
				Name:   "project-id",
				Usage:  "ID of the project to build",
				EnvVar: "BZK_PROJECT_ID",
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
}

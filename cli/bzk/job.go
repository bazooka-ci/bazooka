package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	lib "github.com/haklop/bazooka/commons"

	"github.com/codegangsta/cli"
)

func jobStatus(j lib.JobStatus) string {
	switch j {
	case lib.JOB_SUCCESS:
		return "SUCCESS"
	case lib.JOB_FAILED:
		return "FAILED"
	case lib.JOB_ERRORED:
		return "ERRORED"
	case lib.JOB_RUNNING:
		return "RUNNING"
	default:
		return "-"
	}
}

func fmtTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Local().Format("15:04:05 02/01/2006")
}
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

func listJobsCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list Jobs for the bazooka project",
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
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			res, err := client.ListJobs(c.String("project-id"))
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "JOB ID\tSTARTED\tCOMPLETED\tSTATUS\tPROJECT ID\tORCHESTRATION ID\n")
			for _, item := range res {
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\t%s\t\n", item.ID, fmtTime(item.Started), fmtTime(item.Completed), jobStatus(item.Status), item.ProjectID, item.OrchestrationID)
			}
			w.Flush()
		},
	}
}

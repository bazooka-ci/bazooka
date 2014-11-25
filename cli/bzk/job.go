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
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			res, err := client.StartJob(c.Args()[0], c.Args()[1])
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
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			res, err := client.ListJobs(c.Args()[0])
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

func jobLogCommand() cli.Command {
	return cli.Command{
		Name:  "log",
		Usage: "print the Job's log",
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
				Name:  "job-id",
				Usage: "ID of the job",
			},
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			res, err := client.JobLog(c.String("project-id"), c.String("job-id"))
			if err != nil {
				log.Fatal(err)
			}
			for _, l := range res {
				fmt.Printf("%s [%s] %s\n", l.Time.Format("2006/01/02 15:04:05"), l.Image, l.Message)
			}
		},
	}
}

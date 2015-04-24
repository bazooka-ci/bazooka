package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jawher/mow.cli"

	lib "github.com/bazooka-ci/bazooka/commons"
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

func fmtAuthor(author lib.Person) string {
	if len(author.Email) > 0 {
		if len(author.Name) > 0 {
			return fmt.Sprintf("%s <%s>", author.Name, author.Email)
		}
		return author.Email
	}
	return author.Name
}

func startJobCommand(cmd *cli.Cmd) {
	cmd.Spec = "PROJECT_ID [SCM_REF] [--env...]"

	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "the project id",
	})
	scmRef := cmd.String(cli.StringArg{
		Name:  "SCM_REF",
		Desc:  "the scm ref to build",
		Value: "master",
	})
	envParameters := cmd.Strings(cli.StringsOpt{
		Name: "e env",
		Desc: "define an environment variable for the job",
	})

	cmd.Action = func() {

		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.StartJob(*pid, *scmRef, *envParameters)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "#\tJOB ID\tPROJECT ID\tORCHESTRATION ID\n")
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", res.Number, idExcerpt(res.ID), idExcerpt(res.ProjectID), idExcerpt(res.OrchestrationID))
		w.Flush()
	}

}

func listJobsCommand(cmd *cli.Cmd) {
	cmd.Spec = "[PROJECT_ID]"

	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "the project id",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		var res []lib.Job
		if len(*pid) > 0 {
			res, err = client.ListJobs(*pid)
		} else {
			res, err = client.ListAllJobs()
		}

		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "#\tJOB ID\tSTARTED\tCOMPLETED\tSTATUS\tPROJECT ID\tORCHESTRATION ID\tREFERENCE\tCOMMIT ID\tAUTHOR\tDATE\tMESSAGE\n")
		for _, item := range res {
			fmt.Fprintf(w, "%d\t%s\t%s\t%v\t%v\t%v\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
				item.Number,
				idExcerpt(item.ID),
				fmtTime(item.Started),
				fmtTime(item.Completed),
				jobStatus(item.Status),
				idExcerpt(item.ProjectID),
				idExcerpt(item.OrchestrationID),
				item.SCMMetadata.Reference,
				idExcerpt(item.SCMMetadata.CommitID),
				fmtAuthor(item.SCMMetadata.Author),
				fmtTime(item.SCMMetadata.Date.Time),
				item.SCMMetadata.Message)
		}
		w.Flush()
	}
}

func jobLogCommand(cmd *cli.Cmd) {
	cmd.Spec = "JOB_ID"
	jid := cmd.String(cli.StringArg{
		Name: "JOB_ID",
		Desc: "the job id",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.JobLog(*jid)
		if err != nil {
			log.Fatal(err)
		}
		for _, l := range res {
			fmt.Printf("%s [%s] ", l.Time.Format("2006/01/02 15:04:05"), l.Image)
			switch {
			case len(l.Command) > 0:
				fmt.Printf("[Executing Command] %s\n", l.Command)
			case len(l.Phase) > 0:
				fmt.Printf("[Starting Phase] %s\n", l.Phase)
			case len(l.Level) > 0:
				fmt.Printf("[%s] %s\n", l.Level, l.Message)
			default:
				fmt.Printf("%s\n", l.Message)
			}
		}
	}
}

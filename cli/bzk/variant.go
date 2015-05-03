package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jawher/mow.cli"
)

func listVariantsCommand(cmd *cli.Cmd) {
	cmd.Spec = "JOB_ID"
	jid := cmd.String(cli.StringArg{
		Name: "JOB_ID",
		Desc: "the job id",
	})

	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.Job.Variants(*jid)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)

		fmt.Fprint(w, "NUMBER\tVARIANT ID\tIMAGE\tSTARTED\tCOMPLETED\tSTATUS\tJOB ID\n")
		for _, item := range res {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%v\t%v\t%s\n", item.Number, idExcerpt(item.ID), item.BuildImage, fmtTime(item.Started), fmtTime(item.Completed), jobStatus(item.Status), idExcerpt(item.JobID))
		}
		w.Flush()
	}
}

func variantLogCommand(cmd *cli.Cmd) {
	cmd.Spec = "VARIANT_ID"
	vid := cmd.String(cli.StringArg{
		Name: "VARIANT_ID",
		Desc: "the variant id",
	})
	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.Variant.Log(*vid)
		if err != nil {
			log.Fatal(err)
		}
		for _, l := range res {
			fmt.Printf("%s [%s] %s\n", l.Time.Format("2006/01/02 15:04:05"), l.Image, l.Message)
		}
	}

}

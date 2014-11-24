package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func listVariantsCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list Variant for the bazooka job",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "bazooka-uri",
				Value:  "http://localhost:3000",
				Usage:  "URI for the bazooka server",
				EnvVar: "BZK_URI",
			},
			cli.StringFlag{
				Name:   "project-id",
				Usage:  "ID of the project",
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
			res, err := client.ListVariants(c.String("project-id"), c.String("job-id"))
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)

			fmt.Fprint(w, "NUMBER\tVARIANT ID\tIMAGE\tSTARTED\tCOMPLETED\tSTATUS\tJOB ID\n")
			for _, item := range res {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%v\t%v\t%s\n", item.Number, item.ID, item.BuildImage, fmtTime(item.Started), fmtTime(item.Completed), jobStatus(item.Status), item.JobID)
			}
			w.Flush()
		},
	}
}

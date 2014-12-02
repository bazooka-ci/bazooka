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
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			res, err := client.ListVariants(c.Args()[0])
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)

			fmt.Fprint(w, "NUMBER\tVARIANT ID\tIMAGE\tSTARTED\tCOMPLETED\tSTATUS\tJOB ID\n")
			for _, item := range res {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%v\t%v\t%s\n", item.Number, idExcerpt(item.ID), item.BuildImage, fmtTime(item.Started), fmtTime(item.Completed), jobStatus(item.Status), idExcerpt(item.JobID))
			}
			w.Flush()
		},
	}
}

func variantLogCommand() cli.Command {
	return cli.Command{
		Name:  "log",
		Usage: "print the variant's log",
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
			res, err := client.VariantLog(c.Args()[0])
			if err != nil {
				log.Fatal(err)
			}
			for _, l := range res {
				fmt.Printf("%s [%s] %s\n", l.Time.Format("2006/01/02 15:04:05"), l.Image, l.Message)
			}
		},
	}
}

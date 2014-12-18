package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func addKeyCommand() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "Add SSH Key for the bazooka project",
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
			res, err := client.AddKey(c.Args()[0], c.Args()[1])
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "PROJECT ID\n")
			fmt.Fprintf(w, "%s\n", idExcerpt(res.ProjectID))
			w.Flush()
		},
	}
}

func updateKeyCommand() cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "Update SSH Key for the bazooka project",
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
			res, err := client.UpdateKey(c.Args()[0], c.Args()[1])
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "PROJECT ID\n")
			fmt.Fprintf(w, "%s\n", idExcerpt(res.ProjectID))
			w.Flush()
		},
	}
}

func listKeysCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list Keys for the bazooka project",
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
			res, err := client.ListKeys(c.Args()[0])
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "PROJECT ID\n")
			for _, item := range res {
				fmt.Fprintf(w, "%s\n", idExcerpt(item.ProjectID))
			}
			w.Flush()
		},
	}
}

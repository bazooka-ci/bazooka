package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func listImagesCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list the registered docker images",
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
			res, err := client.ListImages()
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)

			fmt.Fprint(w, "NAME\tIMAGE\tDESCRIPTION\n")
			for _, item := range res {
				fmt.Fprintf(w, "%s\t%s\t%s\n", item.Name, item.Image, item.Description)
			}
			w.Flush()
		},
	}
}

func setImageCommand() cli.Command {
	return cli.Command{
		Name:  "register",
		Usage: "registers a docker image",
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
			if err := client.SetImage(c.Args()[0], c.Args()[1]); err != nil {
				log.Fatal(err)
			}
		},
	}
}

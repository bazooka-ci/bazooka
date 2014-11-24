package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func createProjectCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create a new Project on Bazooka",
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
			res, err := client.CreateProject(c.Args()[0], c.Args()[1], c.Args()[2])
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "PROJECT ID\tNAME\tSCM TYPE\tSCM URI\n")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", res.ID, res.Name, res.ScmType, res.ScmURI)
			w.Flush()
		},
	}
}

func listProjectsCommand() cli.Command {
	return cli.Command{
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
}

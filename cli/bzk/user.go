package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
)

func createUserCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create a new User on Bazooka",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "bazooka-uri",
				Value:  "http://localhost:3000",
				Usage:  "URI for the bazooka server",
				EnvVar: "BZK_URI",
			},
			cli.StringFlag{
				Name:   "password",
				Usage:  "password of the user",
				EnvVar: "BZK_USER_PASSWORD",
			},
		},
		Action: func(c *cli.Context) {
			client, err := NewClient(c.String("bazooka-uri"))
			if err != nil {
				log.Fatal(err)
			}
			password := c.String("password")
			if len(password) == 0 {
				fmt.Printf("Enter user password: ")
				password = string(gopass.GetPasswd())
			}

			res, err := client.CreateUser(c.Args()[0], password)
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "USER ID\tEMAIL\n")
			fmt.Fprintf(w, "%s\t%s\t\n", idExcerpt(res.ID), res.Email)
			w.Flush()
		},
	}
}

func listUsersCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List bazooka users",
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
			res, err := client.ListUsers()
			if err != nil {
				log.Fatal(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
			fmt.Fprint(w, "USER ID\tEMAIL\n")
			for _, item := range res {
				fmt.Fprintf(w, "%s\t%s\t\n", idExcerpt(item.ID), item.Email)
			}
			w.Flush()
		},
	}
}

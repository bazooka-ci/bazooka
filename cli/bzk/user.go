package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jawher/mow.cli"

	"github.com/howeyc/gopass"
)

func createUserCommand(cmd *cli.Cmd) {
	cmd.Spec = "EMAIL [--password]"

	email := cmd.String(cli.StringArg{
		Name: "EMAIL",
		Desc: "The user email",
	})
	password := cmd.String(cli.StringOpt{
		Name:   "p password",
		Desc:   "The user password",
		EnvVar: "BZK_USER_PASSWORD",
	})
	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		if len(*password) == 0 {
			fmt.Printf("Enter user password: ")
			*password = string(gopass.GetPasswd())
		}

		res, err := client.CreateUser(*email, *password)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "USER ID\tEMAIL\n")
		fmt.Fprintf(w, "%s\t%s\t\n", idExcerpt(res.ID), res.Email)
		w.Flush()

	}
}

func listUsersCommand(cmd *cli.Cmd) {
	cmd.Action = func() {
		client, err := NewClient()
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
	}
}

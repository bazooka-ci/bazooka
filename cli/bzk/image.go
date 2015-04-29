package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/jawher/mow.cli"
)

func listImagesCommand(cmd *cli.Cmd) {
	cmd.Action = func() {
		client, err := NewClient()
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
	}
}

func setImageCommand(cmd *cli.Cmd) {
	cmd.Spec = "NAME IMAGE"

	name := cmd.String(cli.StringArg{
		Name: "NAME",
		Desc: "the registration name",
	})
	image := cmd.String(cli.StringArg{
		Name: "IMAGE",
		Desc: "the docker image name",
	})

	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		if err := client.SetImage(*name, *image); err != nil {
			log.Fatal(err)
		}

	}
}

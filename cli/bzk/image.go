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
		res, err := client.Image.List()
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

func getImageCommand(cmd *cli.Cmd) {
	name := cmd.StringArg("IMAGE", "", "The image name")
	cmd.Action = func() {
		client, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}
		image, err := client.Image.Get(*name)
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)

		fmt.Fprint(w, "NAME\tIMAGE\tDESCRIPTION\n")

		fmt.Fprintf(w, "%s\t%s\t%s\n", image.Name, image.Image, image.Description)

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
		if err := client.Image.Set(*name, *image); err != nil {
			log.Fatal(err)
		}

	}
}

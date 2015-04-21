package main

import (
	"fmt"
	"log"

	"github.com/jawher/mow.cli"
)

func encryptData(cmd *cli.Cmd) {
	pid := cmd.String(cli.StringArg{
		Name: "PROJECT_ID",
		Desc: "Project id",
	})
	toEncryptData := cmd.String(cli.StringArg{
		Name: "DATA",
		Desc: "Data to Encrypt",
	})

	cmd.Action = func() {
		client, err := NewClient(checkServerURI(*bzkUri))
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.EncryptData(*pid, *toEncryptData)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Encrypted data: (to add to your .bazooka.yml file)")
		fmt.Printf("%s\n", res)
	}
}

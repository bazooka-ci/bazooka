package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jawher/mow.cli"

	"github.com/howeyc/gopass"
)

func login(cmd *cli.Cmd) {

	email := cmd.String(cli.StringOpt{
		Name:   "email",
		Desc:   "User email",
		EnvVar: "BZK_USER_EMAIL"})
	password := cmd.String(cli.StringOpt{
		Name:   "password",
		Desc:   "User password",
		EnvVar: "BZK_USER_PASSWORD"})

	cmd.Action = func() {
		if len(*email) == 0 {
			*email = interactiveInput("Email")
		}

		if len(*password) == 0 {
			fmt.Printf("Password: ")
			*password = string(gopass.GetPasswd())
		}

		_, err := NewClient(*bzkUri)
		if err != nil {
			log.Fatal(err)
		}

		// TODO check credential. Create a /auth ressource ?

		authConfig := &AuthConfig{
			Username: *email,
			Password: *password,
		}

		saveConfig(authConfig)
	}

}

func interactiveInput(name string) string {
	fmt.Printf("%s: ", name)
	bio := bufio.NewReader(os.Stdin)
	input, isPrefix, err := bio.ReadLine()

	if err != nil {
		log.Fatal("error: ", err)
	}

	if isPrefix {
		log.Fatalf("%s is too long", name)
	}

	return string(input[:])
}

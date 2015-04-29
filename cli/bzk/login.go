package main

import (
	"fmt"
	"log"

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

		_, err := NewClient()
		if err != nil {
			log.Fatal(err)
		}

		// TODO check credential. Create a /auth ressource ?
		config, err := loadConfig()
		if err != nil {
			log.Fatal(fmt.Errorf("Unable to load Bazooka config, reason is: %v\n", err))
		}

		config.Username = *email
		config.Password = *password

		err = saveConfig(config)
		if err != nil {
			log.Fatal(fmt.Errorf("Unable to save Bazooka config, reason is: %v\n", err))
		}
	}

}

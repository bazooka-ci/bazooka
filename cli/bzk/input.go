package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func interactiveInput(name, defaultValue string) string {
	fmt.Printf("%s [%s]:", name, defaultValue)
	bio := bufio.NewReader(os.Stdin)
	input, isPrefix, err := bio.ReadLine()

	if err != nil {
		log.Fatal("error: ", err)
	}

	if isPrefix {
		log.Fatalf("%s is too long", name)
	}

	if len(input) == 0 {
		return defaultValue
	}

	return string(input[:])
}

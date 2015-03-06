package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

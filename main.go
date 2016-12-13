package main

import (
	"log"
	"os"

	"github.com/markbates/grift/cmd"
)

func main() {
	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

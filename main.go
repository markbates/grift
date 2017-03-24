package main

import (
	"os"

	"github.com/markbates/grift/cmd"
)

func main() {
	err := cmd.Run("grift", os.Args[1:])
	if err != nil {
		os.Exit(-1)
	}
}

package cli

const tmpl = `
package main

import (
	"log"
	"os"

	"github.com/markbates/grift/grift"
	_ "{{.Pkg}}"
)

func main() {
	grift.CommandName = "{{.Command}}"
	err := grift.Exec(os.Args[1:], false)
	if err != nil {
		log.Fatal(err)
	}
}
`

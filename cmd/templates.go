package cmd

var modTmpl = `module grifter`

var mainTmpl = `
package main

import _ "{{.GriftsPackagePath}}"
import "os"
import "log"
import "github.com/markbates/grift/grift"
import "path/filepath"

func main() {
	grift.CommandName = "{{.CommandName}}"
	if err := os.Chdir(filepath.Dir("{{.GriftsAbsolutePath}}")); err != nil {
	  log.Fatal(err)
	}
	err := grift.Exec(os.Args[1:], false)
	if err != nil {
		log.Fatal(err)
	}
}`

var initTmpl = `
package grifts

import (
	"fmt"
	"os"
	"strings"

	. "github.com/markbates/grift/grift"
)

var _ = Desc("hello", "Say Hello!")
var _ = Add("hello", func(c *Context) error {
	fmt.Println("Hello World!")
	return nil
})

var _ = Namespace("env", func() {
	Desc("print", "Prints out all of the ENV variables in your environment. Pass in the name of a particular ENV variable to print just that one out. (e.g. grift env:print GOPATH)")
	Add("print", func(c *Context) error {
		if len(c.Args) >= 1 {
			for _, e := range c.Args {
				fmt.Printf("%s=%s\n", e, os.Getenv(e))
			}
		} else {
			for _, e := range os.Environ() {
				pair := strings.Split(e, "=")
				fmt.Printf("%s=%s\n", pair[0], os.Getenv(pair[0]))
			}
		}

		return nil
	})
})
`

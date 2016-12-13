package cmd

var loaderTmpl = `
package grifts

func Load() {}`

var mainTmpl = `
package main

import "{{.GriftsPackagePath}}"
import "os"
import "github.com/markbates/grift/grift"

func main() {
	grifts.Load()
	grift.Exec(os.Args[1:], {{.Verbose}})
}`

var initTmpl = `
package grifts

import (
	"fmt"
	"os"
	"strings"

	. "github.com/markbates/grift/grift"
)

var _ = Add("hello", func(c *Context) error {
	fmt.Println("Hello World!")
	return nil
})


var _ = Desc("env:print", "Prints out all of the ENV variables in your environment. Pass in the name of a particular ENV variable to print just that one out. (e.g. grift env:print GOPATH)")
var _ = Add("env:print", func(c *Context) error {
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
`

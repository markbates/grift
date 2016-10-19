package cmd

var loaderTmpl = `
package tasks

func Load() {}
		`

var mainTmpl = `
package main

import "{{.TasksPackagePath}}"
import "os"
import "log"
import "github.com/markbates/grift/grift"

func main() {
	tasks.Load()
	err := grift.Exec(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
		`

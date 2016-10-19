package main

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

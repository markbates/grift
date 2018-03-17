package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var currentGrift *grifter

func Run(name string, args []string) error {
	defer func() {
		currentGrift.TearDown()
	}()
	if len(args) == 1 {
		switch args[0] {
		case "jim":
			jimTribute()
			return nil
		case "init":
			generateInit()
			return nil
		case "--version", "-v":
			fmt.Println(Version)
			return nil
		case "--help", "-h":
			fmt.Println("grift <task name> [task arguments]")
			return nil
		}
	}

	err := setup(name)
	if err != nil {
		return errors.WithStack(err)
	}

	err = run(args)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func run(args []string) error {
	rargs := []string{"run"}
	// Test for special cases requiring sqlite build tag
	if b, err := ioutil.ReadFile("database.yml"); err == nil {
		if bytes.Contains(b, []byte("sqlite")) {
			rargs = append(rargs, "-tags", "sqlite")
		}
	}
	rargs = append(rargs, exePath)
	rargs = append(rargs, args...)
	if err := grift.RunSource(exec.Command("go", rargs...)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func list() error {
	rargs := []string{"run", exePath, "list"}
	return grift.RunSource(exec.Command("go", rargs...))
}

func setup(name string) error {
	var err error
	currentGrift, err = newGrifter(name)
	if err != nil {
		return errors.WithStack(err)
	}
	err = currentGrift.Setup()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

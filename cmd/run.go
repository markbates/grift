package cmd

import (
	"fmt"
	"os"
	"os/exec"

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
	rargs := []string{"run", exePath}
	rargs = append(rargs, args...)
	runner := exec.Command("go", rargs...)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	err := runner.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func list() error {
	rargs := []string{"run", exePath, "list"}
	runner := exec.Command("go", rargs...)
	runner.Stderr = os.Stderr
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	return runner.Run()
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

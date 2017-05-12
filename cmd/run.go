package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

var currentGrift *grifter

func init() {
	flag.Usage = func() {
		fmt.Printf("Grift Version: %s\n", Version)

		fmt.Print("\nUsage:\n")

		fmt.Println("grift <task name> [task arguments]")

		fmt.Println("\nFlags/Options:")
		flag.PrintDefaults()
	}
}

func Run(name string, args []string) error {
	defer func() {
		currentGrift.TearDown()
	}()
	if len(args) == 2 {
		switch args[1] {
		case "jim":
			jimTribute()
			return nil
		case "init":
			generateInit()
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

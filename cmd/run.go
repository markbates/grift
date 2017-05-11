package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
		return err
	}

	err = run(args)
	if err != nil {
		return err
	}

	return currentGrift.TearDown()
}

func run(args []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.Chdir(currentGrift.BuildPath)
	if err != nil {
		return errors.WithStack(err)
	}
	rargs := []string{"build", "-i", "-o", "grifting"}
	runner := exec.Command("go", rargs...)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	err = runner.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.Chdir(pwd)
	if err != nil {
		return errors.WithStack(err)
	}
	runner = exec.Command(filepath.Join(currentGrift.BuildPath, "grifting"), args...)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	err = runner.Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func list() error {
	rargs := []string{"run", currentGrift.ExePath, "list"}
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
		return err
	}
	err = currentGrift.Setup()
	if err != nil {
		return err
	}
	return currentGrift.Build()
}

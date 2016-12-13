package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var currentGrift *grifter
var verboseFlag bool

func init() {
	flag.Usage = func() {
		fmt.Printf("Grift Version: %s\n", Version)

		fmt.Print("\nUsage:\n")

		fmt.Println("grift [options] <task name> [task arguments]")

		fmt.Println("\nFlags/Options:")
		flag.PrintDefaults()
	}
}

func Run(args []string) error {
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

	flag.BoolVar(&verboseFlag, "v", false, "Print out verbose/debugging information when running a grift")
	flag.Parse()

	err := setup()
	if err != nil {
		return err
	}

	err = run()
	if err != nil {
		return err
	}

	return currentGrift.TearDown()
}

func main() {
	err := Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rargs := []string{"run", currentGrift.ExePath}
	rargs = append(rargs, flag.Args()...)
	runner := exec.Command("go", rargs...)
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	return runner.Run()
}

func list() error {
	rargs := []string{"run", currentGrift.ExePath, "list"}
	runner := exec.Command("go", rargs...)
	runner.Stderr = os.Stderr
	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	return runner.Run()
}

func setup() error {
	var err error
	currentGrift, err = newGrifter()
	if err != nil {
		return err
	}
	currentGrift.Verbose = verboseFlag
	err = currentGrift.Setup()
	if err != nil {
		return err
	}
	return currentGrift.Build()
}

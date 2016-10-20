package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var Version = "0.1.0"
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

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "jim":
			jimTribute()
			return
		case "init":
			generateInit()
			return
		}
	}

	flag.BoolVar(&verboseFlag, "v", false, "Print out verbose/debugging information when running a grift")
	flag.Parse()

	err := setup()
	if err != nil {
		log.Fatal(err)
	}

	err = run()
	if err != nil {
		os.Exit(-1)
		return
	}

	err = currentGrift.TearDown()
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

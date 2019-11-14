package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/markbates/grift/cli"
)

func main() {
	defer func() {
		c := exec.Command("go", "mod", "tidy")
		c.Run()
	}()
	pwd, _ := os.Getwd()
	dir := filepath.Join(pwd, ".grifter")
	defer os.RemoveAll(dir)
	defer func() {
		if err := recover(); err != nil {
			os.RemoveAll(dir)
			log.Fatal(err)
		}
	}()

	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := cli.Main(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

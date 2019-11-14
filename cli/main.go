package cli

import (
	"context"
	"flag"
	"fmt"
)

func Main(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return Help(ctx, args)
	}

	opts := struct {
		version bool
		help    bool
	}{}

	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.BoolVar(&opts.version, "v", false, "display version")
	flags.BoolVar(&opts.help, "h", false, "display help")

	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	if opts.help {
		return Help(ctx, args)
	}

	if opts.version {
		stdout := Stdout(ctx)
		fmt.Fprintln(stdout, Version)
		return nil
	}

	switch args[0] {
	case "init":
		return Init(ctx, args)
	case "jim":
		return Jim(ctx, args)
	}

	return Run(ctx, args)
}

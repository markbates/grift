package cmd

import (
	"context"

	"github.com/markbates/grift/cli"
)

func Run(name string, args []string) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "command", name)
	return cli.Run(ctx, args)
}

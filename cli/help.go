package cli

import (
	"context"
	"fmt"
)

func Help(ctx context.Context, args []string) error {
	stderr := Stderr(ctx)
	fmt.Fprintln(stderr, "grift <task name> [task arguments]")
	return nil
}

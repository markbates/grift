package cli

import (
	"context"
	"io"
	"os"
)

const (
	stdin  = "stdin"
	stdout = "stdout"
	stderr = "stderr"
)

func WithStdin(r io.Reader, ctx context.Context) context.Context {
	return context.WithValue(ctx, stdin, r)
}

func Stdin(ctx context.Context) io.Reader {
	if r, ok := ctx.Value(stdin).(io.Reader); ok {
		return r
	}
	return os.Stdin
}

func Stdout(ctx context.Context) io.Writer {
	if w, ok := ctx.Value(stdout).(io.Writer); ok {
		return w
	}
	return os.Stdout
}

func Stderr(ctx context.Context) io.Writer {
	if w, ok := ctx.Value(stderr).(io.Writer); ok {
		return w
	}
	return os.Stderr
}

func WithStdout(w io.Writer, ctx context.Context) context.Context {
	return context.WithValue(ctx, stdout, w)
}

func WithStderr(w io.Writer, ctx context.Context) context.Context {
	return context.WithValue(ctx, stderr, w)
}

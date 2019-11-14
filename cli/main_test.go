package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Main_Help(t *testing.T) {
	table := [][]string{
		{},
		{"-h"},
	}

	for _, tt := range table {
		t.Run(strings.Join(tt, " "), func(st *testing.T) {
			r := require.New(st)
			bb := &bytes.Buffer{}
			ctx := WithStderr(bb, context.Background())
			r.NoError(Main(ctx, tt))

			r.Contains(bb.String(), "grift <task name> [task arguments]")
		})
	}
}

func Test_Main_Version(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}
	ctx := WithStdout(bb, context.Background())
	r.NoError(Main(ctx, []string{"-v"}))

	r.Equal(Version, strings.TrimSpace(bb.String()))
}

func Test_Main_Task(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}
	ctx := WithStdout(bb, context.Background())

	args := []string{"db:seed", "1", "2", "3"}
	r.NoError(Main(ctx, args))

	act := strings.TrimSpace(bb.String())

	r.Equal("Seeding DB [1 2 3]", act)
}

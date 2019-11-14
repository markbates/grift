package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Run(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	ctx := context.Background()
	ctx = WithStdout(bb, ctx)

	args := []string{"db:seed", "1", "2", "3"}

	err := Run(ctx, args)
	r.NoError(err)

	act := strings.TrimSpace(bb.String())

	r.Equal("Seeding DB [1 2 3]", act)
}

package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Help(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	ctx := WithStderr(bb, context.Background())
	r.NoError(Help(ctx, []string{}))

	r.Contains(bb.String(), "grift <task name> [task arguments]")
}

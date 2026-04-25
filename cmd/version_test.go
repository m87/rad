package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	original := Version
	Version = "dev"
	t.Cleanup(func() { Version = original })

	buf := new(bytes.Buffer)
	versionCmd.SetOut(buf)
	versionCmd.SetErr(buf)
	versionCmd.Run(versionCmd, nil)
	require.Equal(t, "dev\n", buf.String())
}

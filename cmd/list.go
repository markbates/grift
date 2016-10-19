package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
		rargs := []string{"run", currentGrift.ExePath, "list"}
		runner := exec.Command("go", rargs...)
		runner.Stderr = os.Stderr
		runner.Stdin = os.Stdin
		runner.Stdout = os.Stdout
		return runner.Run()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use: "run",
	RunE: func(cmd *cobra.Command, args []string) error {
		rargs := []string{"run", currentGrift.ExePath}
		rargs = append(rargs, os.Args[2:]...)
		runner := exec.Command("go", rargs...)
		runner.Stdin = os.Stdin
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stderr
		return runner.Run()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}

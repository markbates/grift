package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var currentGrift *grifter

var RootCmd = &cobra.Command{
	Use:                "grift",
	DisableFlagParsing: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		currentGrift, err = newGrifter()
		if err != nil {
			return err
		}
		err = currentGrift.Setup()
		if err != nil {
			return err
		}
		return currentGrift.Build()
	},
	PersistentPostRunE: func(c *cobra.Command, args []string) error {
		return currentGrift.TearDown()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

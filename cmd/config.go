package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "print config options",
		Long:  `TODO`,
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
}

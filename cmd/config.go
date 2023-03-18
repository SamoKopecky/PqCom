package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configuration handling",
		Long:  "Used for listing available algorithms and genarating configuration files.",
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
}

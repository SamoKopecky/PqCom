package cmd

import (
	"github.com/spf13/cobra"
)

var (
	log bool

	rootCmd = &cobra.Command{
		Use:   "pqcom",
		Short: "Post quantum communication app",
		Long:  `TODO: Long description`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&log, "log", "l", false, "Enable logging")
}

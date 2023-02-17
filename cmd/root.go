package cmd

import (
	"github.com/spf13/cobra"
)

var (
	log bool

	rootCmd = &cobra.Command{
		Use:   "pqcom",
		Short: "Post quantum communication app",
		Long: `Post quantum communication application for
sending/receiving one time data or chatting`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// TODO: Logging levels
	// rootCmd.PersistentFlags().BoolVarP(&log, "log", "l", false, "Enable logging")
}

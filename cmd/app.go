package cmd

import (
	"github.com/spf13/cobra"
)

var (
	listenPort int
	destPort   int
	destAddr   string

	appCmd = &cobra.Command{
		Use:   "app",
		Short: "Use app mode",
		Long:  `Use the application part of this program`,
	}
)

func init() {
	rootCmd.AddCommand(appCmd)

	appCmd.PersistentFlags().IntVarP(&listenPort, "listen-port", "p", 4040, "Listening port")
	appCmd.PersistentFlags().IntVarP(&destPort, "dest-port", "d", 4040, "Destination port")
	appCmd.PersistentFlags().StringVarP(&destAddr, "dest-addr", "a", "localhost", "Destination address")
}

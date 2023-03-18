package cmd

import (
	"fmt"

	"github.com/SamoKopecky/pqcom/main/app"
	"github.com/spf13/cobra"
)

var (
	connect bool
	listen  bool

	chatCmd = &cobra.Command{
		Use:   "chat",
		Short: "Use chat mode",
		Long:  `Use the application in chat mode where two peers can communicate`,
		Run: func(cmd *cobra.Command, args []string) {
			EnableFileLogging()
			if !connect && !listen {
				fmt.Print("Error: required flag(s) \"connect\" or \"listen\" not set\n\n")
				cmd.Help()
				fmt.Println()
				return
			}
			app.Chat(destAddr, listenPort, destPort, connect)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			LogFile.Close()
		},
	}
)

func init() {
	appCmd.AddCommand(chatCmd)
	chatCmd.Flags().BoolVarP(&connect, "connect", "c", false, "connect to a another user using chap application")
	chatCmd.Flags().BoolVarP(&listen, "listen", "l", false, "listen on a port for another user to connect to you")
	chatCmd.MarkFlagsMutuallyExclusive("connect", "listen")
}

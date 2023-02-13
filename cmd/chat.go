package cmd

import (
	"github.com/SamoKopecky/pqcom/main/app"
	"github.com/spf13/cobra"
)

var (
	chatCmd = &cobra.Command{
		Use:   "chat",
		Short: "Use chat mode",
		Long:  `Use the application in chat mode where two users can communicate`,
		Run: func(cmd *cobra.Command, args []string) {
			app.Chat(destAddr, srcPort, destPort)
		},
	}
)

func init() {
	appCmd.AddCommand(chatCmd)
}

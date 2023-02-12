package cmd

import (
	"github.com/SamoKopecky/pqcom/main/network"
	"github.com/spf13/cobra"
)

var (
	filePath string

	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Use send mode",
		Long: `Use the application in send mode to send data. 
By default the stdin of the app is taken as the source of data.`,
		Run: func(cmd *cobra.Command, args []string) {
			network.Send(destAddr, srcPort, destPort, filePath)
		},
	}
)

func init() {
	appCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVar(&filePath, "file-path", "", "Send data as a file")
}

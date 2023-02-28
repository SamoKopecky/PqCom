package cmd

import (
	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/spf13/cobra"
)

var (
	listenPort int
	destPort   int
	destAddr   string
	configPath string

	appCmd = &cobra.Command{
		Use:   "app",
		Short: "Use app mode",
		Long:  `Use the application part of this program TODO: explain config locations`,
	}
)

func init() {
	rootCmd.AddCommand(appCmd)

	appCmd.PersistentFlags().IntVarP(&listenPort, "listen-port", "p", 4040, "Listening port")
	appCmd.PersistentFlags().IntVarP(&destPort, "dest-port", "d", 4040, "Destination port")
	appCmd.PersistentFlags().StringVarP(&destAddr, "dest-addr", "a", "localhost", "Destination address")
	appCmd.PersistentFlags().StringVar(&configPath, "config", config.DefaultConfigPath, "Config location")
}

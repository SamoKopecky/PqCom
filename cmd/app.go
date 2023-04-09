package cmd

import (
	"fmt"
	"os"

	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/myio"
	"github.com/spf13/cobra"
)

var (
	listenPort int
	destPort   int
	destAddr   string
	configPath string

	defaultConfigPath = fmt.Sprintf("$HOME%c%s%c%s", os.PathSeparator, myio.Config, os.PathSeparator, config.DefaultConfigPath)

	appCmd = &cobra.Command{
		Use:   "app",
		Short: "Use app mode",
		Long: `Use the application part of this program (chatting and file exchange). 
Path to a configuration file can be specified in 3 ways:
1) ENV variable called PQCOM_CONFIG
2) Using the --config flag
3) Default config location at ` + defaultConfigPath,
	}
)

func init() {
	rootCmd.AddCommand(appCmd)

	appCmd.PersistentFlags().IntVarP(&listenPort, "listen-port", "p", 4040, "listening port")
	appCmd.PersistentFlags().IntVarP(&destPort, "dest-port", "d", 4040, "destination port")
	appCmd.PersistentFlags().StringVarP(&destAddr, "dest-addr", "a", "localhost", "destination address")
	appCmd.PersistentFlags().StringVar(&configPath, "config", config.DefaultConfigPath, "config location")
}

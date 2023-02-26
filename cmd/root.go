package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	logOption bool

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
	cobra.OnInitialize(SetLog)
	rootCmd.PersistentFlags().BoolVar(&logOption, "log", false, "Enable logging")
}

func SetLog() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if logOption {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	}
}

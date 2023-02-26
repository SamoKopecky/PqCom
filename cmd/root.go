package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	logOption string

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
	rootCmd.PersistentFlags().StringVar(&logOption, "log", "warning", "Set logging level (trace, debug, info, warning, error, fatal and panic)")
}

func SetLog() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()
	levels := map[string]zerolog.Level{
		"trace":   -1,
		"debug":   0,
		"info":    1,
		"warning": 2,
		"error":   3,
		"fatal":   4,
		"panic":   5,
	}
	keys := make([]string, 0, len(levels))
	for k := range levels {
		keys = append(keys, k)
	}

	if !slices.Contains(keys, logOption) {
		fmt.Println(fmt.Errorf("Uknown log option '%s'", logOption))
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(levels[logOption])
}

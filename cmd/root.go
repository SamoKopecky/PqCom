package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/myio"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var (
	levels = map[string]zerolog.Level{
		"trace":   -1,
		"debug":   0,
		"info":    1,
		"warning": 2,
		"error":   3,
		"fatal":   4,
		"panic":   5,
	}
	LogFile   *os.File
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
	cobra.OnInitialize(SetConfigPath)

	rootCmd.PersistentFlags().StringVar(&logOption, "log", "warning", "set logging level (trace, debug, info, warning, error, fatal and panic)")
}

func SetLog() {
	keys := make([]string, 0, len(levels))
	for k := range levels {
		keys = append(keys, k)
	}

	if !slices.Contains(keys, logOption) {
		fmt.Println(fmt.Errorf("Uknown log option '%s'", logOption))
		os.Exit(1)
	}
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.Level(levels[logOption])).
		With().
		Timestamp().
		Logger()

	// zerolog.SetGlobalLevel(levels[logOption])
}

func EnableFileLogging() {
	dir := myio.HomeSubDir(".local/state")
	filePath := dir + "pqcom_log_" + fmt.Sprintf("%d", time.Now().Unix()) + ".log"
	_, err := os.Stat(dir)
	if err != nil {
		log.Info().Str("dir", filePath).Msg("Creating directory")
		err := os.Mkdir(dir, 0700)
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("Error creating log directory")
		}
	}

	logFile, err := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Error opening log file at ")
	}

	log.Logger = zerolog.New(logFile).
		Level(zerolog.Level(levels[logOption])).
		With().
		Timestamp().
		Logger()
}

func SetConfigPath() {
	config.CmdConfigPath = configPath
}

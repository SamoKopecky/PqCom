package myio

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

const NewFilePerms = 0600

func HomeSubDir(subDir string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msg("Can't get user home dir")
	}
	return strings.Join(
		[]string{home, subDir, ""},
		string(os.PathSeparator))
}

func CreatePath(filePath string) {
	splitPath := strings.Split(filePath, string(os.PathSeparator))
	dirPath := strings.Join(splitPath[:len(splitPath)-1], string(os.PathSeparator))

	_, err := os.Stat(dirPath)
	if err != nil {
		log.Info().Str("path", dirPath).Msg("Creating directory")
		err := os.MkdirAll(dirPath, 0700)
		if err != nil {
			log.Fatal().Str("error", err.Error()).Msg("Error creating log directory")
		}
	}
}

func ContainsFile(file string, dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	for _, e := range entries {
		if file == e.Name() {
			return true, err
		}
	}
	return false, err
}

package config

import (
	"encoding/base64"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Pk []byte
	Sk []byte
}

type Base64Config struct {
	Pk string `json:"public_key"`
	Sk string `json:"private_key"`
}

func ReadConfig() Config {
	configPath := os.Getenv("PQCOM_CONFIG")
	if configPath == "" {
		configPath = "/etc/pqcom_config.json"
	}
	file, err := os.Open(configPath)
	if err != nil {
		log.WithField("error", err).Error("Error opening file")
	}
	defer file.Close()

	base64Config := Base64Config{}
	json.NewDecoder(file).Decode(&base64Config)
	return Config{
		decodeBase64(base64Config.Pk),
		decodeBase64(base64Config.Sk),
	}
}

func decodeBase64(decode string) []byte {
	data, err := base64.StdEncoding.DecodeString(decode)
	if err != nil {
		log.WithField("error", err).Error("Decoding base64")
	}
	return data
}

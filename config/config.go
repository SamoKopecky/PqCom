package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/SamoKopecky/pqcom/main/io"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Kem  string
	Sign string
	Pk   []byte
	Sk   []byte
}

type RawConfig struct {
	Kem  string `json:"kem_alg"`
	Sign string `json:"sign_alg"`
	Pk   string `json:"public_key"`
	Sk   string `json:"private_key"`
}

var CmdConfigPath string

const DefaultConfigPath = "pqcom.json"

func GetConfigPath() string {
	var configPath string
	fmt.Println(CmdConfigPath)
	if CmdConfigPath != DefaultConfigPath {
		configPath = CmdConfigPath
	} else if envConfig := os.Getenv("PQCOM_CONFIG"); envConfig != "" {
		configPath = envConfig
	} else {
		configPath = fmt.Sprintf("%spqcom.json", io.HomeSubDir(".config"))
	}

	_, err := os.Stat(configPath)
	if err != nil {
		log.Fatal().Str("path", configPath).Msg("Config not found")
	}

	log.Info().Str("path", configPath).Msg("Config path set")
	return configPath
}

func ReadConfig() Config {
	file, err := os.Open(GetConfigPath())
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening file")
	}
	defer file.Close()

	rawConfig := RawConfig{}
	json.NewDecoder(file).Decode(&rawConfig)

	// TODO: use generics for getAll{Kems | Signs} and use it here
	if crypto.IsValidAlg(rawConfig.Kem, crypto.GetKemNames) {
		log.Info().Str("algorithm", rawConfig.Kem).Msg("Using key encapsulation to exchange keys")
	} else {
		log.Fatal().Str("algorithm", rawConfig.Kem).Msg("Unkown key encapsulation in config")
	}

	if crypto.IsValidAlg(rawConfig.Sign, crypto.GetSignNames) {
		log.Info().Str("algorithm", rawConfig.Sign).Msg("Using signature to secure key exchange")
	} else {
		log.Fatal().Str("algorithm", rawConfig.Sign).Msg("Unkown signature")
	}

	decodedPk := decodeBase64(rawConfig.Pk)
	decodedSk := decodeBase64(rawConfig.Sk)
	sign := crypto.GetSign(rawConfig.Sign).F
	if rawConfig.Pk == "" || len(decodedPk) != sign.PkLen() {
		log.Fatal().Msg("Incorrect length of the configured public key")
	}

	if rawConfig.Sk == "" || len(decodedSk) != sign.SkLen() {
		log.Fatal().Msg("Incorrect length of the configured private key")
	}

	return Config{
		rawConfig.Kem,
		rawConfig.Sign,
		decodedPk,
		decodedSk,
	}
}

func decodeBase64(decode string) []byte {
	data, err := base64.StdEncoding.DecodeString(decode)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Decoding base64")
	}
	return data
}

func GenerateConfig(kem, sign string) {
	cPk, cSk := crypto.GenerateKeys(sign)
	sPk, sSk := crypto.GenerateKeys(sign)
	clientConfig := RawConfig{kem, sign, sPk, cSk}
	serverConfig := RawConfig{kem, sign, cPk, sSk}
	writeConfig(clientConfig, "./pqcom_client_example.json")
	writeConfig(serverConfig, "./pqcom_server_example.json")
}

func writeConfig(rawConfig RawConfig, name string) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening file")
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	encoder.Encode(rawConfig)
	file.Close()
}

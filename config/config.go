package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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

func ReadConfig() Config {
	configPath := os.Getenv("PQCOM_CONFIG")
	if configPath == "" {
		configPath = "/etc/pqcom_config.json"
	}
	log.WithField("path", configPath).Info("Loaded config")
	file, err := os.Open(configPath)
	if err != nil {
		log.WithField("error", err).Fatal("Error opening file")
	}
	defer file.Close()

	rawConfig := RawConfig{}
	json.NewDecoder(file).Decode(&rawConfig)

	// TODO: use generics for getAll{Kems | Signs} and use it here
	if crypto.IsValidAlg(rawConfig.Kem, crypto.GetAllKems) {
		log.WithField("algorithm", rawConfig.Kem).Info("Using key encapsulation to exchange keys")
	} else {
		log.WithField("algorithm", rawConfig.Kem).Fatal("Unkown key encapsulation in config")
	}

	if crypto.IsValidAlg(rawConfig.Sign, crypto.GetAllSigns) {
		log.WithField("algorithm", rawConfig.Sign).Info("Using signature to secure key exchange")
	} else {
		log.WithField("algorithm", rawConfig.Sign).Fatal("Unkown signature")
	}

	decodedPk := decodeBase64(rawConfig.Pk)
	decodedSk := decodeBase64(rawConfig.Sk)
	sign := crypto.GetSign(rawConfig.Sign).Functions
	fmt.Println(sign.PkLen())
	if rawConfig.Pk == "" || len(decodedPk) != sign.PkLen() {
		log.Fatal("Incorrect length of the configured public key")
	}

	if rawConfig.Sk == "" || len(decodedSk) != sign.SkLen() {
		log.Fatal("Incorrect length of the configured private key")
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
		log.WithField("error", err).Fatal("Decoding base64")
	}
	return data
}

func GenerateConfig(kem, sign string) {
	pk, sk := crypto.GenerateKeys(sign)
	config := RawConfig{
		kem,
		sign,
		pk,
		sk,
	}
	file, err := os.OpenFile("./pqcom_config_example.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.WithField("error", err).Fatal("Error opening file")
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	encoder.Encode(config)
	file.Close()
}
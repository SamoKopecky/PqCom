package config

import (
	"encoding/base64"
	"encoding/json"
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
	log.WithField("algorithm", rawConfig.Kem).Info("Using key encapsulation to exchange keys")
	log.WithField("algorithm", rawConfig.Sign).Info("Using signature to secure key exchange")
	return Config{
		rawConfig.Kem,
		rawConfig.Sign,
		decodeBase64(rawConfig.Pk),
		decodeBase64(rawConfig.Sk),
	}
}

func decodeBase64(decode string) []byte {
	data, err := base64.StdEncoding.DecodeString(decode)
	if err != nil {
		log.WithField("error", err).Error("Decoding base64")
	}
	return data
}

func GenerateConfig(kem, sign string) {
	pk, sk := GenerateKeys(sign)
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

func GenerateKeys(sign string) (pkStr string, skStr string) {
	signFncs := crypto.GetSign(sign).Functions
	pk, sk := signFncs.KeyGen()
	pkStr = base64.StdEncoding.EncodeToString(pk)
	skStr = base64.StdEncoding.EncodeToString(sk)
	return
}

package crypto

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/dilithium"
	"github.com/sirupsen/logrus"
)

func generateKeys() (pkStr string, skStr string) {
	pk, sk := dilithium.KeyGen()
	pkStr = base64.StdEncoding.EncodeToString(pk)
	skStr = base64.StdEncoding.EncodeToString(sk)
	return
}

func WriteKeys() {
	pk, sk := generateKeys()
	config := config.Base64Config{Pk: pk, Sk: sk}
	file, err := os.OpenFile("./pqcom_config_example.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.WithField("error", err).Error("Error opening file")
	}
	json.NewEncoder(file).Encode(config)
}

package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	log "github.com/sirupsen/logrus"
)

const NONCE_LEN = 12

func GenerateNonce() []byte {
	nonce := make([]byte, NONCE_LEN)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.WithField("error", err).Fatal("Generating nonce")
	}
	return nonce
}

func getRowIndex[T any](to map[string]T, key string) int {
	keys := make([]string, 0, len(to))
	for k := range to {
		keys = append(keys, k)
	}
	for i, k := range keys {
		if k == key {
			return i
		}
	}
	return -1
}

func GenerateKeys(sign string) (pkStr string, skStr string) {
	signFncs := GetSign(sign).Functions
	pk, sk := signFncs.KeyGen()
	pkStr = base64.StdEncoding.EncodeToString(pk)
	skStr = base64.StdEncoding.EncodeToString(sk)
	return
}

func IsValidAlg(option string, getAll func() []string) bool {
	for _, key := range getAll() {
		if key == option {
			return true
		}
	}
	return false
}

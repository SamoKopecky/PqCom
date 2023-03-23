package crypto

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/rs/zerolog/log"
)

type AesCipher struct {
	gcm cipher.AEAD
}

func (aesCipher *AesCipher) Create(key []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Creating aes cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Creating AES GCM struct")
	}
	aesCipher.gcm = gcm
}

func (aesCipher *AesCipher) Encrypt(data []byte) ([]byte, []byte) {
	nonce := GenerateNonce()
	return aesCipher.gcm.Seal(nil, nonce, data, nil), nonce
}

func (aesCipher *AesCipher) Decrypt(data, nonce []byte) []byte {
	plaintext, err := aesCipher.gcm.Open(nil, nonce, data, nil)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Decrypting")
	}
	return plaintext
}

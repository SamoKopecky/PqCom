package crypto

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/rs/zerolog/log"
)

type AesCipher struct {
	gcm   cipher.AEAD
	nonce []byte
}

func (aesCipher *AesCipher) Create(key []byte, nonce []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Creating aes cipher")
	}

	aesCipher.nonce = nonce

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Creating AES GCM struct")
	}
	aesCipher.gcm = gcm
}

func (aesCipher *AesCipher) Encrypt(data []byte) []byte {
	return aesCipher.gcm.Seal(nil, aesCipher.nonce, data, nil)
}

func (aesCipher *AesCipher) Decrypt(data []byte) []byte {
	plaintext, err := aesCipher.gcm.Open(nil, aesCipher.nonce, data, nil)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Decrypting")
	}
	return plaintext
}

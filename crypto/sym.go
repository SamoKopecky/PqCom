package crypto

import (
	"crypto/aes"
	"crypto/cipher"

	log "github.com/sirupsen/logrus"
)

type AesCipher struct {
	gcm   cipher.AEAD
	nonce []byte
}

func (aesCipher *AesCipher) Create(key []byte, nonce []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.WithField("error", err).Fatal("Creating aes cipher")
	}

	aesCipher.nonce = nonce

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.WithField("error", err).Fatal("Creating AES GCM struct")
	}
	aesCipher.gcm = gcm
}

func (aesCipher *AesCipher) Encrypt(data []byte) []byte {
	return aesCipher.gcm.Seal(nil, aesCipher.nonce, data, nil)
}

func (aesCipher *AesCipher) Decrypt(data []byte) []byte {
	plaintext, err := aesCipher.gcm.Open(nil, aesCipher.nonce, data, nil)
	if err != nil {
		log.WithField("error", err).Fatal("Decrypting")
	}
	return plaintext
}

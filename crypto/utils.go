package crypto

import (
	"crypto/rand"
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

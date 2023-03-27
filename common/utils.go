package common

import "golang.org/x/crypto/sha3"

func PMod(i, m int) (o int) {
	return (i%m + m) % m
}

func Kdf(input []byte, len int) (output []byte) {
	output = make([]byte, len)
	sha3.ShakeSum256(output, input)
	return
}

package dilithium

import "golang.org/x/crypto/sha3"

func shake256(input []byte, len int) (output []byte) {
	output = make([]byte, len)
	sha3.ShakeSum256(output, input)
	return
}

func shake128(input []byte, len int) (output []byte) {
	output = make([]byte, len)
	sha3.ShakeSum128(output, input)
	return
}
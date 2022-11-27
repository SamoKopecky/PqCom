package kyber

import "golang.org/x/crypto/sha3"

func xof(input []byte, x byte, y byte, len int) (output []byte) {
	output = make([]byte, len)
	input = append(input, x)
	input = append(input, y)
	sha3.ShakeSum256(output, input)
	return
}

func hash32(input []byte) (hashedBytes []byte) {
	hashedBytes = make([]byte, 32)
	sha := sha3.Sum256(input)
	for i := 0; i < 32; i++ {
		hashedBytes[i] = sha[i]
	}
	return
}

func hash64(input []byte) (first []byte, second []byte) {
	hashedBytes := sha3.Sum512(input)
	first = hashedBytes[:32]
	second = hashedBytes[32:]
	return
}

func kdf(input []byte, len int) (output []byte) {
	output = make([]byte, len)
	sha3.ShakeSum256(output, input)
	return
}

func prf(input []byte, localN byte, len int) (output []byte) {
	output = make([]byte, len)
	input = append(input, localN)
	sha3.ShakeSum256(output, input)
	return
}

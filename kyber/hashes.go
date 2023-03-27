package kyber

import "golang.org/x/crypto/sha3"

func xof(input []byte, x byte, y byte, len int) (output []byte) {
	output = make([]byte, len)
	input = append(input, x)
	input = append(input, y)
	sha3.ShakeSum256(output, input)
	return
}

func hash32(input []byte) []byte {
	temp := sha3.Sum256(input)
	return temp[:]
}

func hash64(input []byte) (first, second []byte) {
	output := sha3.Sum512(input)
	first = output[:32]
	second = output[32:]
	return
}

func prf(input []byte, localN byte, len int) (output []byte) {
	output = make([]byte, len)
	input = append(input, localN)
	sha3.ShakeSum256(output, input)
	return
}

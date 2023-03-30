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

func BytesToBits(bytes []byte) (bits []byte) {
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			bits = append(bits, ExtractBit(int(bytes[i]), j))
		}
	}
	return
}

func ExtractBit(from int, power int) (bit byte) {
	bit = byte(from & (1 << power) >> power)
	return
}

package common

import "golang.org/x/crypto/sha3"

const n = 256

func PMod(i, m int) (o int) {
	return (i%m + m) % m
}

func Kdf(input []byte, len int) (output []byte) {
	output = make([]byte, len)
	sha3.ShakeSum256(output, input)
	return
}

func BytesToBits(bytes []byte) (bits []byte) {
	var i, I, j int
	bits = make([]byte, len(bytes)*8)
	for i = 0; i < len(bytes); i++ {
		I = i * 8
		for j = 0; j < 8; j++ {
			bits[I+j] = ExtractBit(int(bytes[i]), j)
		}
	}
	return
}
func PolyToBits(poly []int, coefSize int) (bits []byte) {
	var i, I, j int
	bits = make([]byte, n*coefSize)
	for i = 0; i < n; i++ {
		I = i * coefSize
		for j = 0; j < coefSize; j++ {
			bits[I+j] = ExtractBit((poly[i]), j)
		}
	}
	return
}

func ExtractBit(from int, power int) (bit byte) {
	bit = byte(from & (1 << power) >> power)
	return
}

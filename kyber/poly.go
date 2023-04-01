package kyber

import (
	"math"

	"github.com/SamoKopecky/pqcom/main/common"
)

func (kyb *Kyber) encode(poly []int, coefSize int) (bytes []byte) {
	bits := common.PolyToBits(poly, coefSize)

	var number, mask byte
	var i, I, j int

	polyBits := n * coefSize
	polyBytes := polyBits / 8

	bytes = make([]byte, polyBytes)
	for i = 0; i < polyBits; i += 8 {
		number = 0
		mask = 1
		for j = 0; j < 8; j++ {
			number += bits[j+i] * mask
			mask <<= 1
		}
		bytes[I] = number
		I++
	}
	return
}

func (kyb *Kyber) decode(bytes []byte, coefSize int) (poly []int) {
	bits := common.BytesToBits(bytes)

	poly = make([]int, n)

	var coef, i, I, j, mask int
	for i = 0; i < n*coefSize; i += coefSize {
		coef = 0
		mask = 1
		for j = 0; j < coefSize; j++ {
			coef += int(bits[i+j]) * mask
			mask <<= 1
		}
		poly[I] = coef
		I++
	}
	return
}

func (kyb *Kyber) compress(input []int, d int) (compressed []int) {
	var value int
	modulo := 1 << d
	temp := float64(modulo) / float64(q)
	for _, v := range input {
		value = int(math.Round(temp * float64(v)))
		compressed = append(compressed, common.PMod(value, modulo))
	}
	return
}

func (kyb *Kyber) decompress(input []int, d int) (decompressed []int) {
	divisor := float64(int(1 << d))
	for _, v := range input {
		decompressed = append(decompressed, int(math.Round(q/divisor*float64(v))))
	}
	return
}

func (kyb *Kyber) add(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] + g[i])
	}
	return
}

func (kyb *Kyber) sub(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] - g[i])
	}
	return
}

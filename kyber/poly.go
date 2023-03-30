package kyber

import (
	"math"

	"github.com/SamoKopecky/pqcom/main/common"
)

func (kyb *Kyber) encode(poly []int, l int) (bytes []byte) {
	bits := make([]byte, l*n)
	for i := 0; i < n; i++ {
		for j := 0; j < l; j++ {
			bits[j+i*l] = common.ExtractBit(poly[i], j)
		}
	}
	var encoded byte
	var k int
	bytes = make([]byte, l*n/8)
	for i := 0; i < l*n; i += 8 {
		for j := 0; j < 8; j++ {
			encoded += (bits[j+i]) * (1 << j)
		}
		bytes[k] = encoded
		k++
		encoded = 0
	}
	return
}

func (kyb *Kyber) decode(bytes []byte, l int) (poly []int) {
	bits := common.BytesToBits(bytes)
	poly = make([]int, n)
	var fi int
	for i := 0; i < n; i++ {
		fi = 0
		for j := 0; j < l; j++ {
			fi += int(bits[i*l+j]) * (1 << j)
		}
		poly[i] = fi
	}
	return
}

func (kyb *Kyber) compress(input []int, d int) (compressed []int) {
	var value int
	modulo := 1 << d
	moduloFloat := float64(modulo)
	temp := moduloFloat / float64(q)
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

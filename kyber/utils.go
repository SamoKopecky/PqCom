package kyber

import (
	"crypto/rand"
	"math"
)

func cbd(bytes []byte, eta int) (poly []int) {
	bits := bytesToBits(bytes)
	for i := 0; i < 256; i++ {
		a := 0
		b := 0
		for j := 0; j < eta; j++ {
			a += int(bits[2*i*eta+j])
		}
		for j := 0; j < eta; j++ {
			b += int(bits[2*i*eta+eta+j])
		}
		poly = append(poly, a-b)
	}
	return
}

func parse(bytes []byte) (ntt_poly []int) {
	j, i := 0, 0
	ntt_poly = make([]int, N)
	for j < N {
		reduced := modPlus(int(bytes[i+1]), 16)
		d1 := int(bytes[i]) + 256*reduced
		d2 := int(math.Floor(float64(bytes[i+1]))) + int((16 * bytes[i+2]))
		if d1 < Q {
			ntt_poly[j] = d1
			j++
		}
		if d2 < Q && j < N {
			ntt_poly[j] = d2
			j++
		}
		i = i + 3
	}
	return
}

func bytesToBits(bytes []byte) (bits []byte) {
	for i := 0; i < len(bytes)*8; i++ {
		bits = append(bits, (bytes[i/8]/byte(math.Pow(2, float64(i%8))))%2)
	}
	return
}

func randBytes(size int) (randBytes []byte) {
	randBytes = make([]byte, size)
	rand.Read(randBytes)
	return
}

func BytesEqual(a, b []byte) (equal bool) {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func modPlus(a int, modulo int) (b int) {
	b = a % modulo
	if a < 0 {
		b += modulo
	}
	return
}

func randPoly(r []byte, localN *byte, eta int) []int {
	return cbd(prf(r, *localN, 64*eta), eta)
}

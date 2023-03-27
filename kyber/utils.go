package kyber

import (
	"crypto/rand"
)

func (kyb *Kyber) cbd(bytes []byte, eta int) (poly []int) {
	var a, b int
	bits := kyb.bytesToBits(bytes)
	for i := 0; i < n; i++ {
		a = 0
		b = 0
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

func (kyb *Kyber) parse(bytes []byte) (ntt_poly []int) {
	var j, i, d1, d2 int
	ntt_poly = make([]int, n)
	for j < n {
		d1 = int(bytes[i]) + n*kyb.pMod(int(bytes[i+1]), 16)
		d2 = int(bytes[i+1]/16) + int((16 * bytes[i+2]))
		if d1 < q {
			ntt_poly[j] = d1
			j++
		}
		if d2 < q && j < n {
			ntt_poly[j] = d2
			j++
		}
		i += 3
	}
	return
}

func (kyb *Kyber) bytesToBits(bytes []byte) (bits []byte) {
	for i := 0; i < len(bytes)*8; i++ {
		bits = append(bits, (bytes[i/8]/(1<<(i&0x7)))&0x1)
	}
	return
}

func (kyb *Kyber) randBytes(size int) (randBytes []byte) {
	randBytes = make([]byte, size)
	rand.Read(randBytes)
	return
}

func (kyb *Kyber) BytesEqual(a, b []byte) (equal bool) {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (kyb *Kyber) pMod(i, m int) (o int) {
	return (i%m + m) % m
}

func (kyb *Kyber) randPoly(r []byte, localN byte, eta int) []int {
	return kyb.cbd(prf(r, localN, 64*eta), eta)
}

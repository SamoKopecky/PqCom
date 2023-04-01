package kyber

import (
	"crypto/rand"

	"github.com/SamoKopecky/pqcom/main/common"
)

func (kyb *Kyber) cbd(bytes []byte, eta int) (poly []int) {
	var a, b, Eta int
	bits := common.BytesToBits(bytes)
	poly = make([]int, n)

	for i := 0; i < n; i++ {
		Eta = 2 * i * eta
		a = 0
		b = 0
		for j := 0; j < eta; j++ {
			a += int(bits[Eta+j])
		}
		for j := 0; j < eta; j++ {
			b += int(bits[Eta+eta+j])
		}
		poly[i] = a - b
	}
	return
}

func (kyb *Kyber) parse(bytes []byte) (ntt_poly []int) {
	var j, i, d1, d2 int
	ntt_poly = make([]int, n)
	for j < n {
		d1 = int(bytes[i]) + n*common.PMod(int(bytes[i+1]), 16)
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

func (kyb *Kyber) randBytes(size int) (randBytes []byte) {
	randBytes = make([]byte, size)
	rand.Read(randBytes)
	return
}

func (kyb *Kyber) randPoly(r []byte, localN byte, eta int) []int {
	return kyb.cbd(kyb.prf(r, localN, 64*eta), eta)
}

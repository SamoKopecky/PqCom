package dilithium

import (
	"crypto/rand"
	"encoding/binary"
	"math"

	"golang.org/x/crypto/sha3"
)

func (dil *Dilithium) modP(a, mod int) (o int) {
	o = ((a % mod) + mod) % mod
	return
}

func (dil *Dilithium) modPM(a, mod int) int {
	first := int(math.Abs(float64(((a % mod) - mod) % mod)))
	if first <= mod/2 {
		return int(-first)
	}
	return int(dil.modP(a, mod))
}

func (dil *Dilithium) powerToRound(r int) (int, int) {
	r = dil.modP(r, q)
	r0 := dil.modPM(r, 1<<d)
	return (r - r0) / (1 << d), r0
}

func (dil *Dilithium) decompose(r, alpha int) (r_1, r_0 int) {
	r = dil.modP(r, q)
	r_0 = dil.modPM(r, alpha)
	if r-r_0 == q-1 {
		r_1 = 0
		r_0 -= 1
	} else {
		r_1 = int((r - r_0) / alpha)
	}
	return
}

func (dil *Dilithium) highBits(r, alpha int) (r_1 int) {
	r_1, _ = dil.decompose(r, alpha)
	return
}

func (dil *Dilithium) lowBits(r, alpha int) (r_0 int) {
	_, r_0 = dil.decompose(r, alpha)
	return
}

func (dil *Dilithium) makeHint(z, r, alpha int) bool {
	return dil.highBits(r, alpha) != dil.highBits(r+z, alpha)
}

func (dil *Dilithium) useHint(h bool, r, alpha int) int {
	m := int((q - 1) / alpha)
	r1, r0 := dil.decompose(r, alpha)
	if h && r0 > 0 {
		return dil.modP(r1+1, m)
	} else if h && r0 <= 0 {
		return dil.modP(r1-1, m)
	}
	return r1
}

func (dil *Dilithium) sampleInBall(c_wave []byte) (c []int) {
	c = make([]int, 256)
	shake := sha3.NewShake256()
	shake.Write(c_wave)
	o := make([]byte, 8)
	shake.Read(o)

	bits := dil.bytesToBits(o)[:dil.tau]
	for i := 256 - dil.tau; i < 256; i++ {
		j_byte := make([]byte, 1)
		j := byte(255)
		for j > byte(i) {
			shake.Read(j_byte)
			j = j_byte[0]
		}
		s := bits[i-(256-dil.tau)]
		c[i] = c[j]
		c[j] = int(1 - int8(2*s))
	}
	return
}

func (dil *Dilithium) bytesToBits(bytes []byte) (bits []byte) {
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			bits = append(bits, dil.extractBit(int(bytes[i]), j))
		}
	}
	return
}

func (dil *Dilithium) polyToBits(poly []int, l int) (bits []byte) {
	for i := 0; i < 256; i++ {
		for j := 0; j < l; j++ {
			bits = append(bits, dil.extractBit((poly[i]), j))
		}
	}
	return
}

func (dil *Dilithium) extractBit(from int, power int) (bit byte) {
	bit = byte(from & (1 << power) >> power)
	return
}

func (dil *Dilithium) genRand(bits int) (xofOutput []byte) {
	len := bits / 8
	xofOutput = make([]byte, len)
	randBytes := make([]byte, 256)
	rand.Read(randBytes)
	xofOutput = dil.shake256(randBytes, len)
	return
}

func (dil *Dilithium) expandS(ro_dash []byte) (vectors [][]int) {
	for i := 0; i < dil.l+dil.k; i++ {
		poly := []int{}
		shake := sha3.NewShake256()
		i_bytes := make([]byte, 2)
		i_bytes[0] = byte(i)
		i_bytes[1] = byte(i >> 8)
		shake.Write(ro_dash)
		shake.Write(i_bytes)

		for len(poly) < n {
			o := make([]byte, 1)
			shake.Read(o)

			two_ints := [...]uint8{uint8(o[0]) >> 4, uint8(o[0]) & 0xF}

			if dil.eta == 2 {
				for _, v := range two_ints {
					if v >= 15 {
						continue
					}
					poly = append(poly, dil.eta-(int(v)%5))
				}
			}
			if dil.eta == 4 {
				for _, v := range two_ints {
					v := int(v)
					if v >= 2*dil.eta+1 {
						continue
					}
					poly = append(poly, dil.eta-v)
				}
			}
		}
		vectors = append(vectors, poly[:n])
	}
	return
}

func (dil *Dilithium) expandMask(ro_dash []byte, kappa int) (y [][]int) {
	for i := 0; i < dil.l; i++ {
		poly := make([]int, n)
		shake := sha3.NewShake256()
		shake.Write(ro_dash)
		sum := kappa + i
		bytes := make([]byte, 2)
		bytes[0] = byte(sum)
		bytes[1] = byte(sum >> 8)
		shake.Write(bytes)

		for j := 0; j < n; j++ {
			o := make([]byte, 4)
			shake.Read(o)
			o[0] = 0
			o[1] = o[1] & 2
			poly[j] = dil.gammaOne - int(binary.BigEndian.Uint32(o))
		}
		y = append(y, poly)
	}
	return
}

func (dil *Dilithium) BytesEqual(a, b []byte) (equal bool) {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

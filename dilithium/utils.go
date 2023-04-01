package dilithium

import (
	"crypto/rand"

	"github.com/SamoKopecky/pqcom/main/common"
	"golang.org/x/crypto/sha3"
)

func (dil *Dilithium) modPM(a, mod int) int {
	mMod := ((a % mod) - mod) % mod
	mMod = dil.abs(mMod)
	if 2*mMod <= mod {
		return int(-mMod)
	}

	return mod - mMod
}

func (dil *Dilithium) powerToRound(r int) (int, int) {
	r = common.PMod(r, q)
	r0 := dil.modPM(r, 1<<d)
	return (r - r0) / (1 << d), r0
}

func (dil *Dilithium) decompose(r, alpha int) (r_1, r_0 int) {
	r = common.PMod(r, q)
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
		return common.PMod(r1+1, m)
	} else if h && r0 <= 0 {
		return common.PMod(r1-1, m)
	}
	return r1
}

func (dil *Dilithium) sampleInBall(c_wave []byte) (c []int) {
	var j_byte []byte

	c = make([]int, n)
	shake := sha3.NewShake256()
	shake.Write(c_wave)
	o := make([]byte, 8)
	shake.Read(o)
	j := byte(255)

	bits := common.BytesToBits(o)[:dil.tau]
	for i := n - dil.tau; i < n; i++ {
		j_byte = make([]byte, 1)
		for j > byte(i) {
			shake.Read(j_byte)
			j = j_byte[0]
		}
		s := bits[i-(n-dil.tau)]
		c[i] = c[j]
		c[j] = int(1 - int8(2*s))
	}
	return
}

func (dil *Dilithium) genRand(bits int) (xofOutput []byte) {
	len := bits / 8
	xofOutput = make([]byte, len)
	randBytes := make([]byte, n)
	rand.Read(randBytes)
	xofOutput = common.Kdf(randBytes, len)
	return
}

func (dil *Dilithium) expandS(ro_dash []byte) (vectors [][]int) {
	vectors = make([][]int, dil.l+dil.k)
	i_bytes := make([]byte, 2)
	o := make([]byte, 1)
	data := make([]byte, n)
	var shake sha3.ShakeHash
	var i, j int
	var two_ints [2]byte

	for i = 0; i < dil.l+dil.k; i++ {
		vectors[i] = make([]int, n)
		j = 0
		shake = sha3.NewShake256()
		i_bytes[0] = byte(i)
		i_bytes[1] = byte(uint16(i) >> 8)
		shake.Write(ro_dash)
		shake.Write(i_bytes)
		shake.Read(data)

		for j < n {
			o = data[j : j+1]
			two_ints = [2]byte{byte(o[0]) >> 4, byte(o[0]) & 0xF}

			if dil.eta == 2 {
				for _, v := range two_ints {
					if j >= n {
						break
					}
					if v >= 15 {
						o = data[j : j+1]
						continue
					}
					vectors[i][j] = dil.eta - (int(v) % 5)
				}
			} else {
				// Eta == 4
				for _, v := range two_ints {
					if j >= n {
						break
					}
					if v >= 9 {
						o = data[j : j+1]
						continue
					}
					vectors[i][j] = dil.eta - int(v)
				}
			}
			j++
		}
	}
	return
}

func (dil *Dilithium) expandMask(ro_dash []byte, kappa int) (vec [][]int) {
	bytes := make([]byte, 2)
	data := make([]byte, 3*n)
	o := make([]byte, 3)
	vec = make([][]int, dil.l)
	var sum int
	var shake sha3.ShakeHash

	var restOfBitsAndOp = byte(0x3)
	if dil.gammaOneExp == 19 {
		restOfBitsAndOp = byte(0xF)
	}

	for i := 0; i < dil.l; i++ {
		vec[i] = make([]int, n)
		shake = sha3.NewShake256()
		shake.Write(ro_dash)
		sum = kappa + i
		bytes[0] = byte(sum)
		bytes[1] = byte(uint16(sum) >> 8)
		shake.Write(bytes)
		shake.Read(data)

		for j := 0; j < n; j++ {
			o = data[3*j : 3*(j+1)]
			o[2] &= restOfBitsAndOp
			vec[i][j] = dil.gammaOne - dil.littleEndian(o)
		}
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

func (dil *Dilithium) littleEndian(bytes []byte) int {
	return int(uint32(bytes[0]) | (uint32(bytes[1]) << 8) | (uint32(bytes[2]) << 16))
}

func (dil *Dilithium) abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

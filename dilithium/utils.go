package dilithium

import (
	"crypto/rand"
	"encoding/binary"
	"math"

	"golang.org/x/crypto/sha3"
)

func modP(a, mod int) int {
	return ((a % mod) + mod) % mod
}

func PowerToMod(r, d int) (int, int) {
	r = modP(r, Q)
	two_power_d := int(math.Pow(2.0, float64(d)))
	r0 := r % two_power_d
	return (r - r0) / two_power_d, r0
}

func Decompose(r, alpha int) (r1, r0 int) {
	r = modP(r, Q)
	r0 = r % alpha
	if r-r0 == Q-1 {
		r1 = 0
		r0 -= 1
	} else {
		r1 = int((r - r0) / alpha)
	}
	return
}

func HighBits(r, alpha int) (r1 int) {
	r1, _ = Decompose(r, alpha)
	return
}

func LowBits(r, alpha int) (r0 int) {
	_, r0 = Decompose(r, alpha)
	return
}

func MakeHint(z, r, alpha int) bool {
	return HighBits(r, alpha) != HighBits(r+z, alpha)
}

func UseHint(h bool, r, alpha int) int {
	m := int((Q - 1) / alpha)
	r1, r0 := Decompose(r, alpha)
	if h && r0 > 0 {
		return modP(r1+1, m)
	} else if h && r0 <= 0 {
		return modP(r1-1, m)
	}
	return r1
}

// func SampleInBall(ro []byte) {
// 	c := make([]byte, 256)
// 	for i := 256-TAU; i < 256; i++ {
// 		j := xof(ro, i)
// 		s := xof(ro, 1)
// 		c[i] = c[j]

// 	}
// }

func genRand(bits int) (xofOutput []byte) {
	len := bits / 8
	xofOutput = make([]byte, len)
	randBytes := make([]byte, 256)
	rand.Read(randBytes)
	xofOutput = shake256(randBytes, len)
	return
}

func expandS(ro_dash []byte) (vectors [][]int) {
	for i := 0; i < L+K; i++ {
		poly := []int{}
		shake := sha3.NewShake256()
		i_bytes := []byte{0x00, byte(i)}
		to_shake := append(ro_dash, i_bytes...)
		shake.Write(to_shake)
		for len(poly) < N {
			o := make([]byte, 1)
			shake.Read(o)

			int_output := uint8(o[0])
			two_ints := [...]uint8{int_output >> 4, int_output & 0xF}

			if Eta == 2 {
				for _, v := range two_ints {
					if v >= 15 {
						continue
					}
					poly = append(poly, Eta-modP(int(v), 5))
				}
			}
			// TODO: eta = 4
		}
		vectors = append(vectors, poly[:N])
	}
	return
}

func expandA(ro []byte) (mat [][][]int) {
	for i := 0; i < K; i++ {
		row := [][]int{}
		for j := 0; j < L; j++ {
			poly := []int{}
			shake := sha3.NewShake128()
			shake.Write(ro)
			i_and_j := [2]byte{byte(i), byte(j)}
			shake.Write(i_and_j[:])
			for len(poly) < N {
				// TODO: make this so it uses little endian
				// Right it works correctly but is written confusing
				o_3 := make([]byte, 3)
				shake.Read(o_3)
				o_3[0] = o_3[0] & 0x7F
				zero := [1]byte{0}
				o_4 := append(zero[:], o_3...)
				parsed := int(binary.BigEndian.Uint32(o_4))
				if parsed > Q-1 {
					continue
				}
				poly = append(poly, parsed)
			}
			row = append(row, poly)
		}
		mat = append(mat, row)
	}
	return
}

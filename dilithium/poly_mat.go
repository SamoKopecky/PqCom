package dilithium

import (
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

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

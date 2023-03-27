package dilithium

import (
	"golang.org/x/crypto/sha3"
)

func (dil *Dilithium) expandA(ro []byte) (mat [][][]int) {
	var poly []int
	var row [][]int
	var shake sha3.ShakeHash
	var i, j byte

	for i = 0; i < byte(dil.k); i++ {
		row = [][]int{}
		for j = 0; j < byte(dil.l); j++ {
			poly = []int{}
			shake = sha3.NewShake128()
			shake.Write(ro)
			shake.Write([]byte{i})
			shake.Write([]byte{j})

			for len(poly) < n {
				o_3 := make([]byte, 3)
				shake.Read(o_3)
				o_3[2] &= 0x7F
				num := dil.littleEndian(o_3)
				if num > q-1 {
					continue
				}
				poly = append(poly, num)
			}
			row = append(row, poly)
		}
		mat = append(mat, row)
	}
	return
}

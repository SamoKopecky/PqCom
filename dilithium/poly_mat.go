package dilithium

import (
	"golang.org/x/crypto/sha3"
)

func (dil *Dilithium) expandA(ro []byte) (mat [][][]int) {
	var shake sha3.ShakeHash
	var i, j byte
	var k int
	var o_3 []byte
	var data = make([]byte, 3*n)
	mat = make([][][]int, dil.k)

	for i = 0; i < byte(dil.k); i++ {
		mat[i] = make([][]int, dil.l)
		for j = 0; j < byte(dil.l); j++ {
			mat[i][j] = make([]int, n)
			k = 0
			shake = sha3.NewShake128()
			shake.Write(ro)
			shake.Write([]byte{i, j})
			shake.Read(data)

			for k < n {
				o_3 = data[3*k : 3*(k+1)]
				o_3[2] &= 0x7F
				num := dil.littleEndian(o_3)
				if num >= q {
					shake.Read(data)
					continue
				}
				mat[i][j][k] = num
				k++
			}
		}
	}
	return
}

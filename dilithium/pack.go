package dilithium

import "github.com/SamoKopecky/pqcom/main/common"

func (dil *Dilithium) bitPackPolyVec(a [][]int, coefSize int) (bytes []byte) {
	var number byte
	var i, I, j, J, k int
	polyBits := n * coefSize
	polyBytes := polyBits / 8
	vecs := len(a)

	vecBits := make([][]byte, vecs)
	for i = 0; i < vecs; i++ {
		vecBits[i] = common.PolyToBits(a[i], coefSize)
	}

	bits := make([]byte, polyBits)
	bytes = make([]byte, polyBytes*vecs)

	for i = 0; i < vecs; i++ {
		bits = vecBits[i]
		I = i * polyBytes
		J = 0
		for j = 0; j < polyBits; j += 8 {
			number = 0
			for k = 0; k < 8; k++ {
				number += bits[j+k] * (1 << k)
			}
			bytes[I+J] = number
			J++
		}
	}
	return
}

func (dil *Dilithium) bitUnpackPolyVec(bytes []byte, coefSize int) (polyVec [][]int) {
	bits := common.BytesToBits(bytes)

	polyC := len(bits) / coefSize / n
	polyVec = make([][]int, polyC)

	var coef, i, I, j, J, k int
	vecBits := n * coefSize

	for i = 0; i < polyC; i++ {
		polyVec[i] = make([]int, n)
		I = i * vecBits
		J = 0
		for j = 0; j < vecBits; j += coefSize {
			coef = 0
			for k = 0; k < coefSize; k++ {
				coef += int(bits[I+j+k]) * (1 << k)
			}
			polyVec[i][J] = coef
			J++
		}
	}
	return
}

func (dil *Dilithium) bitPackAlteredPolyVec(a [][]int, alter, size int) (o []byte) {
	b := make([][]int, len(a))

	for i := 0; i < len(b); i++ {
		b[i] = make([]int, n)
		for j := 0; j < len(b[0]); j++ {
			b[i][j] = alter - a[i][j]
		}
	}
	o = dil.bitPackPolyVec(b, size)
	return
}

func (dil *Dilithium) bitUnpackAlteredPolyVec(bytes []byte, alter, size int) (o [][]int) {
	a := dil.bitUnpackPolyVec(bytes, size)
	poly := make([]int, n)

	for i := 0; i < len(a); i++ {
		poly = make([]int, n)
		for j := 0; j < len(a[0]); j++ {
			poly[j] = (alter - a[i][j])
		}
		o = append(o, poly)
	}

	return
}

func (dil *Dilithium) bitPackHint(h [][]byte) (o []byte) {
	ones_len := 0
	lengths := make([]byte, len(h))
	var row_positions []byte
	var row_len int
	for i := 0; i < len(h); i++ {
		row_positions = []byte{}
		for j := 0; j < len(h[0]); j++ {
			if h[i][j] == 1 {
				row_positions = append(row_positions, byte(j))
			}
		}
		row_len = len(row_positions)
		lengths[i] = byte(row_len)
		ones_len += row_len
		o = append(o, row_positions...)
	}
	padding := make([]byte, dil.omega-ones_len)
	o = append(o, padding...)
	o = append(o, lengths...)

	return
}

func (dil *Dilithium) bitUnpackHint(bytes []byte) (h [][]byte) {
	lengths := bytes[len(bytes)-dil.k:]
	start := 0
	end := int(lengths[0]) - 1
	for i := 0; i < dil.k; i++ {
		row := make([]byte, n)
		for j := start; j <= end; j++ {
			row[bytes[j]] = 1
		}
		h = append(h, row)
		if i == dil.k-1 {
			continue
		}
		start += int(lengths[i])
		end += int(lengths[i+1])
	}
	return
}

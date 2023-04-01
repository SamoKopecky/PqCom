package dilithium

import "github.com/SamoKopecky/pqcom/main/common"

func (dil *Dilithium) bitPackPolyVec(a [][]int, coefSize int) (bytes []byte) {
	var number, mask byte
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
			mask = 1
			for k = 0; k < 8; k++ {
				number += bits[j+k] * mask
				mask <<= 1
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

	var coef, i, I, j, J, k, mask int
	vecBits := n * coefSize

	for i = 0; i < polyC; i++ {
		polyVec[i] = make([]int, n)
		I = i * vecBits
		J = 0
		for j = 0; j < vecBits; j += coefSize {
			coef = 0
			mask = 1
			for k = 0; k < coefSize; k++ {
				coef += int(bits[I+j+k]) * mask
				mask <<= 1
			}
			polyVec[i][J] = coef
			J++
		}
	}
	return
}

func (dil *Dilithium) bitPackAlteredPolyVec(input [][]int, alter, size int) []byte {
	temp := make([][]int, len(input))

	for i := 0; i < len(temp); i++ {
		temp[i] = make([]int, n)
		for j := 0; j < len(temp[0]); j++ {
			temp[i][j] = alter - input[i][j]
		}
	}
	return dil.bitPackPolyVec(temp, size)
}

func (dil *Dilithium) bitUnpackAlteredPolyVec(bytes []byte, alter, size int) (output [][]int) {
	temp := dil.bitUnpackPolyVec(bytes, size)
	output = make([][]int, len(temp))

	for i := 0; i < len(temp); i++ {
		output[i] = make([]int, n)
		for j := 0; j < len(temp[0]); j++ {
			output[i][j] = (alter - temp[i][j])
		}
	}
	return
}

func (dil *Dilithium) bitPackHint(hints [][]byte) (output []byte) {
	ones_len := 0
	lengths := make([]byte, len(hints))
	var row_positions []byte
	var row_len int
	for i := 0; i < len(hints); i++ {
		row_positions = []byte{}
		for j := 0; j < len(hints[0]); j++ {
			if hints[i][j] == 1 {
				row_positions = append(row_positions, byte(j))
			}
		}
		row_len = len(row_positions)
		lengths[i] = byte(row_len)
		ones_len += row_len
		output = append(output, row_positions...)
	}
	padding := make([]byte, dil.omega-ones_len)
	output = append(output, padding...)
	output = append(output, lengths...)

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

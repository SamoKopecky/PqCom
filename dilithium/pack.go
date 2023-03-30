package dilithium

import "github.com/SamoKopecky/pqcom/main/common"

func (dil *Dilithium) bitPackPolyVec(a [][]int, size int) (o []byte) {
	var number byte
	var i, j, k, l int
	onePolyLen := n * size / 8

	bitsRows := make([][]byte, len(a))
	for i = 0; i < len(a); i++ {
		bitsRows[i] = dil.polyToBits(a[i], size)
	}
	bits := make([]byte, n*size)
	o = make([]byte, onePolyLen*len(a))
	for i = 0; i < len(bitsRows); i++ {
		bits = bitsRows[i]
		l = 0
		for j = 0; j < len(bits); j += 8 {
			number = 0
			for k = 0; k < 8; k++ {
				number += bits[j+k] * (1 << k)
			}
			o[(i*onePolyLen)+l] = number
			l++
		}
	}
	return
}

func (dil *Dilithium) bitUnpackPolyVec(bytes []byte, size int) (o [][]int) {
	bits := common.BytesToBits(bytes)
	o = make([][]int, len(bits)/size/256)
	var number, l, m, i, j, k int

	for i = 0; i < len(bits); i += n * size {
		o[m] = make([]int, n)
		l = 0
		for j = 0; j < n*size; j += size {
			number = 0
			for k = 0; k < size; k++ {
				number += int(bits[i+j+k]) * (1 << k)
			}
			o[m][l] = number
			l++
		}
		m++
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

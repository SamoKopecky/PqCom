package dilithium

func bitPackPolyVec(a [][]int, size int) (o []byte) {
	bits := []byte{}
	for i := 0; i < len(a); i++ {
		bits = append(bits, polyToBits(a[i], size)...)
	}
	for i := 0; i < len(bits); i += 8 {
		number := byte(0)
		for j := 0; j < 8; j++ {
			number += bits[i+j] * 1 << j
		}
		o = append(o, number)
	}
	return
}

func bitUnpackPolyVec(bytes []byte, size int) (o [][]int) {
	bits := bytesToBits(bytes)
	for i := 0; i < len(bits); i += N * size {
		poly := []int{}
		for j := 0; j < N*size; j += size {
			number := 0
			for k := 0; k < size; k++ {
				number += int(bits[i+j+k]) * 1 << k
			}
			poly = append(poly, number)
		}
		o = append(o, poly)
	}
	return
}

func bitPackAlteredPolyVec(a [][]int, alter, size int) (o []byte) {
	b := make([][]int, len(a))

	for i := 0; i < len(b); i++ {
		b[i] = make([]int, N)
		for j := 0; j < len(b[0]); j++ {
			b[i][j] = alter - a[i][j]
		}
	}
	o = bitPackPolyVec(b, size)
	return
}

func bitUnpackAlteredPolyVec(bytes []byte, alter, size int) (o [][]int) {
	a := bitUnpackPolyVec(bytes, size)

	for i := 0; i < len(a); i++ {
		poly := make([]int, N)
		for j := 0; j < len(a[0]); j++ {
			poly[j] = (alter - a[i][j])
		}
		o = append(o, poly)
	}

	return
}

func bitPackHint(h [][]byte) (o []byte) {
	ones_len := 0
	lengths := make([]byte, len(h))
	for i := 0; i < len(h); i++ {
		row_positions := []byte{}
		for j := 0; j < len(h[0]); j++ {
			if h[i][j] == 1 {
				row_positions = append(row_positions, byte(j))
			}
		}
		row_len := len(row_positions)
		lengths[i] = byte(row_len)
		ones_len += row_len
		o = append(o, row_positions...)
	}
	padding := make([]byte, Omega-ones_len)
	o = append(o, padding...)
	o = append(o, lengths...)

	return
}

func bitUnpackHint(bytes []byte) (h [][]byte) {
	lengths := bytes[len(bytes)-4:]
	start := 0
	end := int(lengths[0]) - 1
	for i := 0; i < K; i++ {
		row := make([]byte, N)
		for j := start; j <= end; j++ {
			row[bytes[j]] = 1
		}
		h = append(h, row)
		if i == K-1 {
			continue
		}
		start += int(lengths[i])
		end += int(lengths[i+1])
	}
	return
}
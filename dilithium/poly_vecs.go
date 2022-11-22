package dilithium

import (
	"math"
)

func mulPolyVec(f, g [][]int) (h []int) {
	h = make([]int, N)
	for i := 0; i < K; i++ {
		h = polyAdd(polyMul(f[i], g[i]), h)
	}
	return
}

func addPolyVec(f, g [][]int) (h [][]int) {
	h = make([][]int, K)
	for i := 0; i < K; i++ {
		h[i] = polyAdd(f[i], g[i])
	}
	return
}

func subPolyVec(f, g [][]int) (h [][]int) {
	h = make([][]int, K)
	for i := 0; i < K; i++ {
		h[i] = polySub(f[i], g[i])
	}
	return
}

func inversePolyVec(f [][]int) (h [][]int) {
	h = make([][]int, K)
	for i := 0; i < K; i++ {
		h_row := make([]int, N)
		for j := 0; j < len(f[0]); j++ {
			h_row[j] = -f[i][j]
		}
		h[i] = h_row
	}
	return
}

func nttPolyVec(polyVec [][]int) (polyVecCopy [][]int) {
	polyVecCopy = make([][]int, len(polyVec))
	copy(polyVecCopy, polyVec)

	for i := 0; i < K; i++ {
		polyVecCopy[i] = ntt(polyVecCopy[i])
	}
	return
}

func invNttPolyVec(polyVec [][]int) (polyVecCopy [][]int) {
	polyVecCopy = make([][]int, len(polyVec))
	copy(polyVecCopy, polyVec)

	for i := 0; i < K; i++ {
		polyVecCopy[i] = invNtt(polyVecCopy[i])
	}
	return
}

func scalePolyVecByPoly(a [][]int, b []int) (o [][]int) {
	for i := 0; i < len(a); i++ {
		o = append(o, polyMul(a[i], b))
	}
	return
}

func scalePolyVecByInt(a [][]int, b int) (o [][]int) {
	for i := 0; i < len(a); i++ {
		row := make([]int, N)
		for j := 0; j < N; j++ {
			row[j] = (a[i][j] * b) % Q
		}
		o = append(o, row)
	}
	return
}

func powerToModPolyVec(r [][]int, d int) (r_1, r_2 [][]int) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]int, N)
		r2_row := make([]int, N)
		for j := 0; j < len(r[0]); j++ {
			r_r_1, r_r_2 := powerToRound(r[i][j], d)
			r1_row[j] = r_r_1
			r2_row[j] = r_r_2
		}
		r_1 = append(r_1, r1_row)
		r_2 = append(r_2, r2_row)
	}
	return
}

func makeHintPolyVec(z, r [][]int, alpha int) (r_1 [][]byte) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]byte, N)
		for j := 0; j < len(r[0]); j++ {
			if makeHint(z[i][j], r[i][j], alpha) {
				r1_row[j] = 1
			}
		}
		r_1 = append(r_1, r1_row)
	}
	return
}

func useHintPolyVec(h [][]byte, r [][]int, alpha int) (o [][]int) {
	h_bool := false
	for i := 0; i < len(r); i++ {
		o_row := make([]int, N)
		for j := 0; j < len(r[0]); j++ {
			h_bool = h[i][j] != 0
			o_row[j] = useHint(h_bool, r[i][j], alpha)
		}
		o = append(o, o_row)
	}
	return
}

func highBitsPolyVec(r [][]int, alpha int) (r_1 [][]int) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]int, N)
		for j := 0; j < len(r[0]); j++ {
			r1_row[j] = highBits(r[i][j], alpha)
		}
		r_1 = append(r_1, r1_row)
	}
	return
}

func lowBitsPolyVec(r [][]int, alpha int) (r_1 [][]int) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]int, N)
		for j := 0; j < len(r[0]); j++ {
			r1_row[j] = lowBits(r[i][j], alpha)
		}
		r_1 = append(r_1, r1_row)
	}
	return
}

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

func modPPolyVec(a [][]int) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < N; j++ {
			a[i][j] = modP(a[i][j], Q)
		}
	}
	return
}

func reducePolyVec(a [][]int) (b [][]int) {
	for i := 0; i < len(a); i++ {
		row := make([]int, N)
		for j := 0; j < N; j++ {
			row[j] = modPM(a[i][j], Q)
		}
		b = append(b, row)
	}
	return
}

func checkNormPolyVec(a [][]int, bound int) bool {
	max := 0
	abs := 0
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			test := int(math.Abs(float64(a[i][j])))
			abs = test % Q
			if abs > max {
				max = abs
			}
		}
	}
	if max >= bound {
		return true
	}
	return false
}

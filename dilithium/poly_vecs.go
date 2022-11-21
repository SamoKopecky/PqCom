package dilithium

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

func powerToModPolyVec(r [][]int, d int) (r_1, r_2 [][]int) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]int, N)
		r2_row := make([]int, N)
		for j := 0; j < len(r[0]); j++ {
			r_r_1, r_r_2 := powerToMod(r[i][j], d)
			r1_row[j] = r_r_1
			r2_row[j] = r_r_2
		}
		r_1 = append(r_1, r1_row)
		r_2 = append(r_2, r2_row)
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
			poly[j] = alter - a[i][j]
		}
		o = append(o, poly)
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

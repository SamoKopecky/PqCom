package dilithium

func (dil *Dilithium) mulPolyVec(f, g [][]int) (h []int) {
	h = make([]int, n)
	for i := 0; i < dil.l; i++ {
		h = dil.polyAdd(dil.polyMul(f[i], g[i]), h)
	}
	return
}

func (dil *Dilithium) addPolyVec(f, g [][]int) (h [][]int) {
	h = make([][]int, len(f))
	for i := 0; i < len(f); i++ {
		h[i] = dil.polyAdd(f[i], g[i])
	}
	return
}

func (dil *Dilithium) subPolyVec(f, g [][]int) (h [][]int) {
	h = make([][]int, len(f))
	for i := 0; i < len(f); i++ {
		h[i] = dil.polySub(f[i], g[i])
	}
	return
}

func (dil *Dilithium) inversePolyVec(f [][]int) (h [][]int) {
	h = make([][]int, dil.k)
	for i := 0; i < dil.k; i++ {
		h_row := make([]int, n)
		for j := 0; j < len(f[0]); j++ {
			h_row[j] = -f[i][j]
		}
		h[i] = h_row
	}
	return
}

func (dil *Dilithium) nttPolyVec(polyVec [][]int) (polyVecCopy [][]int) {
	polyVecCopy = make([][]int, len(polyVec))
	copy(polyVecCopy, polyVec)

	for i := 0; i < len(polyVec); i++ {
		polyVecCopy[i] = dil.ntt(polyVecCopy[i])
	}
	return
}

func (dil *Dilithium) invNttPolyVec(polyVec [][]int) (polyVecCopy [][]int) {
	polyVecCopy = make([][]int, len(polyVec))
	copy(polyVecCopy, polyVec)

	for i := 0; i < len(polyVec); i++ {
		polyVecCopy[i] = dil.invNtt(polyVecCopy[i])
	}
	return
}

func (dil *Dilithium) scalePolyVecByPoly(a [][]int, b []int) (o [][]int) {
	aLen := len(a)
	o = make([][]int, aLen)
	for i := 0; i < aLen; i++ {
		o[i] = dil.polyMul(a[i], b)
	}
	return
}

func (dil *Dilithium) scalePolyVecByInt(a [][]int, b int) (o [][]int) {
	aLen := len(a)
	o = make([][]int, aLen)
	for i := 0; i < aLen; i++ {
		o[i] = make([]int, n)
		for j := 0; j < n; j++ {
			o[i][j] = (a[i][j] * b) % q
		}
	}
	return
}

func (dil *Dilithium) powerToModPolyVec(r [][]int) (r_1, r_2 [][]int) {
	rLen := len(r)
	rRowLen := len(r[0])
	r_1 = make([][]int, rLen)
	r_2 = make([][]int, rLen)
	for i := 0; i < rLen; i++ {
		r_1[i] = make([]int, n)
		r_2[i] = make([]int, n)
		for j := 0; j < rRowLen; j++ {
			r_r_1, r_r_2 := dil.powerToRound(r[i][j])
			r_1[i][j] = r_r_1
			r_2[i][j] = r_r_2
		}
	}
	return
}

func (dil *Dilithium) makeHintPolyVec(z, r [][]int, alpha int) (r_1 [][]byte) {
	rLen := len(r)
	rRowLen := len(r[0])
	r_1 = make([][]byte, rLen)
	for i := 0; i < rLen; i++ {
		r_1[i] = make([]byte, n)
		for j := 0; j < rRowLen; j++ {
			if dil.makeHint(z[i][j], r[i][j], alpha) {
				r_1[i][j] = 1
			}
		}
	}
	return
}

func (dil *Dilithium) useHintPolyVec(h [][]byte, r [][]int, alpha int) (o [][]int) {
	h_bool := false
	rLen := len(r)
	rRowLen := len(r[0])
	o = make([][]int, rLen)
	for i := 0; i < rLen; i++ {
		o[i] = make([]int, n)
		for j := 0; j < rRowLen; j++ {
			h_bool = h[i][j] != 0
			o[i][j] = dil.useHint(h_bool, r[i][j], alpha)
		}
	}
	return
}

func (dil *Dilithium) highBitsPolyVec(r [][]int, alpha int) (r_1 [][]int) {
	rLen := len(r)
	rRowLen := len(r[0])
	r_1 = make([][]int, rLen)
	for i := 0; i < rLen; i++ {
		r_1[i] = make([]int, n)
		for j := 0; j < rRowLen; j++ {
			r_1[i][j] = dil.highBits(r[i][j], alpha)
		}
	}
	return
}

func (dil *Dilithium) lowBitsPolyVec(r [][]int, alpha int) (r_1 [][]int) {
	rLen := len(r)
	rRowLen := len(r[0])
	r_1 = make([][]int, rLen)
	for i := 0; i < rLen; i++ {
		r_1[i] = make([]int, n)
		for j := 0; j < rRowLen; j++ {
			r_1[i][j] = dil.lowBits(r[i][j], alpha)
		}
	}
	return
}

func (dil *Dilithium) modPMPolyVec(a [][]int, mod int) (b [][]int) {
	aLen := len(a)
	b = make([][]int, aLen)
	for i := 0; i < aLen; i++ {
		b[i] = make([]int, n)
		for j := 0; j < n; j++ {
			b[i][j] = dil.PMmod(a[i][j], mod)
		}
	}
	return
}

func (dil *Dilithium) checkNormPolyVec(a [][]int, bound int) bool {
	max := 0
	abs := 0
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			abs = dil.abs(a[i][j])
			if abs > max {
				max = abs
			}
		}
	}
	return max >= bound
}

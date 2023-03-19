package dilithium

import (
	"math"
)

func (dil *dilithium) mulPolyVec(f, g [][]int) (h []int) {
	h = make([]int, n)
	for i := 0; i < dil.k; i++ {
		h = dil.polyAdd(dil.polyMul(f[i], g[i]), h)
	}
	return
}

func (dil *dilithium) addPolyVec(f, g [][]int) (h [][]int) {
	h = make([][]int, dil.k)
	for i := 0; i < dil.k; i++ {
		h[i] = dil.polyAdd(f[i], g[i])
	}
	return
}

func (dil *dilithium) subPolyVec(f, g [][]int) (h [][]int) {
	h = make([][]int, dil.k)
	for i := 0; i < dil.k; i++ {
		h[i] = dil.polySub(f[i], g[i])
	}
	return
}

func (dil *dilithium) inversePolyVec(f [][]int) (h [][]int) {
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

func (dil *dilithium) nttPolyVec(polyVec [][]int) (polyVecCopy [][]int) {
	polyVecCopy = make([][]int, len(polyVec))
	copy(polyVecCopy, polyVec)

	for i := 0; i < dil.k; i++ {
		polyVecCopy[i] = dil.ntt(polyVecCopy[i])
	}
	return
}

func (dil *dilithium) invNttPolyVec(polyVec [][]int) (polyVecCopy [][]int) {
	polyVecCopy = make([][]int, len(polyVec))
	copy(polyVecCopy, polyVec)

	for i := 0; i < dil.k; i++ {
		polyVecCopy[i] = dil.invNtt(polyVecCopy[i])
	}
	return
}

func (dil *dilithium) scalePolyVecByPoly(a [][]int, b []int) (o [][]int) {
	for i := 0; i < len(a); i++ {
		o = append(o, dil.polyMul(a[i], b))
	}
	return
}

func (dil *dilithium) scalePolyVecByInt(a [][]int, b int) (o [][]int) {
	for i := 0; i < len(a); i++ {
		row := make([]int, n)
		for j := 0; j < n; j++ {
			row[j] = (a[i][j] * b) % q
		}
		o = append(o, row)
	}
	return
}

func (dil *dilithium) powerToModPolyVec(r [][]int) (r_1, r_2 [][]int) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]int, n)
		r2_row := make([]int, n)
		for j := 0; j < len(r[0]); j++ {
			r_r_1, r_r_2 := dil.powerToRound(r[i][j])
			r1_row[j] = r_r_1
			r2_row[j] = r_r_2
		}
		r_1 = append(r_1, r1_row)
		r_2 = append(r_2, r2_row)
	}
	return
}

func (dil *dilithium) makeHintPolyVec(z, r [][]int, alpha int) (r_1 [][]byte) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]byte, n)
		for j := 0; j < len(r[0]); j++ {
			if dil.makeHint(z[i][j], r[i][j], alpha) {
				r1_row[j] = 1
			}
		}
		r_1 = append(r_1, r1_row)
	}
	return
}

func (dil *dilithium) useHintPolyVec(h [][]byte, r [][]int, alpha int) (o [][]int) {
	h_bool := false
	for i := 0; i < len(r); i++ {
		o_row := make([]int, n)
		for j := 0; j < len(r[0]); j++ {
			h_bool = h[i][j] != 0
			o_row[j] = dil.useHint(h_bool, r[i][j], alpha)
		}
		o = append(o, o_row)
	}
	return
}

func (dil *dilithium) highBitsPolyVec(r [][]int, alpha int) (r_1 [][]int) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]int, n)
		for j := 0; j < len(r[0]); j++ {
			r1_row[j] = dil.highBits(r[i][j], alpha)
		}
		r_1 = append(r_1, r1_row)
	}
	return
}

func (dil *dilithium) lowBitsPolyVec(r [][]int, alpha int) (r_1 [][]int) {
	for i := 0; i < len(r); i++ {
		r1_row := make([]int, n)
		for j := 0; j < len(r[0]); j++ {
			r1_row[j] = dil.lowBits(r[i][j], alpha)
		}
		r_1 = append(r_1, r1_row)
	}
	return
}

func (dil *dilithium) modPPolyVec(a [][]int) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < n; j++ {
			a[i][j] = dil.modP(a[i][j], q)
		}
	}
	return
}

func (dil *dilithium) modPMPolyVec(a [][]int, mod int) (b [][]int) {
	for i := 0; i < len(a); i++ {
		row := make([]int, n)
		for j := 0; j < n; j++ {
			row[j] = dil.modPM(a[i][j], mod)
		}
		b = append(b, row)
	}
	return
}

func (dil *dilithium) checkNormPolyVec(a [][]int, bound int) bool {
	max := 0
	abs := 0
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			abs = int(math.Abs(float64(a[i][j])))
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

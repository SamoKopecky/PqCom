package kyber

import "github.com/SamoKopecky/pqcom/main/common"

func (kyb *Kyber) mulVec(f, g [][]int) (h []int) {
	h = make([]int, n)
	for i := 0; i < kyb.k; i++ {
		h = kyb.add(kyb.pointWiseMulVec(f[i], g[i]), h)
	}
	return
}

func (kyb *Kyber) addVec(f, g [][]int) (h [][]int) {
	h = make([][]int, kyb.k)
	for i := 0; i < kyb.k; i++ {
		h[i] = kyb.add(f[i], g[i])
	}
	return
}

func (kyb *Kyber) modPVec(a [][]int) {
	for i := 0; i < kyb.k; i++ {
		for j := 0; j < n; j++ {
			a[i][j] = common.PMod(a[i][j], q)
		}
	}
}

func (kyb *Kyber) randPolyVec(r []byte, localN *byte, eta int) (vector [][]int) {
	vector = make([][]int, kyb.k)
	for i := 0; i < kyb.k; i++ {
		vector[i] = kyb.randPoly(r, *localN, eta)
		*localN++
	}
	return
}

func (kyb *Kyber) decodePolyVec(bytes []byte, coefSize int) (polyVec [][]int) {
	polyVec = make([][]int, kyb.k)
	interval := coefSize * n / 8
	var I int

	for i := 0; i < interval*kyb.k; i += interval {
		polyVec[I] = kyb.decode(bytes[i:i+interval], coefSize)
		I++
	}
	return
}

func (kyb *Kyber) encodePolyVec(polyVec [][]int, coefSize int) (bytes []byte) {
	for i := 0; i < kyb.k; i++ {
		// Append here is fine
		bytes = append(bytes, kyb.encode(polyVec[i], coefSize)...)
	}
	return
}

func (kyb *Kyber) nttPolyVec(polyVec [][]int) {
	for i := 0; i < kyb.k; i++ {
		kyb.ntt(polyVec[i])
	}
}

func (kyb *Kyber) invNttPolyVec(polyVec [][]int) {
	for i := 0; i < kyb.k; i++ {
		kyb.invNtt(polyVec[i])
	}
}

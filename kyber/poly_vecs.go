package kyber

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
			a[i][j] = kyb.pMod(a[i][j], q)
		}
	}
}

func (kyb *Kyber) randPolyVec(r []byte, localN *byte, eta int) (vector [][]int) {
	vector = [][]int{}
	for i := 0; i < kyb.k; i++ {
		vector = append(vector, kyb.randPoly(r, *localN, eta))
		*localN++
	}
	return
}

func (kyb *Kyber) decodePolyVec(bytes []byte, l int) (polyVec [][]int) {
	polyVec = make([][]int, kyb.k)
	interval := l * n / 8
	j := 0

	for i := 0; i < interval*kyb.k; i += interval {
		polyVec[j] = kyb.decode(bytes[i:i+interval], l)
		j++
	}
	return
}

func (kyb *Kyber) encodePolyVec(polyVec [][]int, l int) (bytes []byte) {
	for i := 0; i < kyb.k; i++ {
		bytes = append(bytes, kyb.encode(polyVec[i], l)...)
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

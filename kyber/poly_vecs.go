package kyber

func (kyb *kyber) mulPolyVec(f, g [][]int) (h []int) {
	h = make([]int, n)
	for i := 0; i < kyb.k; i++ {
		h = kyb.polyAdd(kyb.polyMul(f[i], g[i]), h)
	}
	return
}

func (kyb *kyber) addPolyVec(f, g [][]int) (h [][]int) {
	h = make([][]int, kyb.k)
	for i := 0; i < kyb.k; i++ {
		h[i] = kyb.polyAdd(f[i], g[i])
	}
	return
}

func (kyb *kyber) modPlusPolyVec(a [][]int) {
	for i := 0; i < kyb.k; i++ {
		for j := 0; j < n; j++ {
			a[i][j] = kyb.modPlus(a[i][j], q)
		}
	}
	return
}

func (kyb *kyber) randPolyVec(r []byte, localN *byte, eta int) (vector [][]int) {
	vector = [][]int{}
	for i := 0; i < kyb.k; i++ {
		vector = append(vector, kyb.randPoly(r, localN, eta))
		*localN += 1
	}
	return
}

func (kyb *kyber) decodePolyVec(bytes []byte, l int) (polyVec [][]int) {
	polyVec = make([][]int, kyb.k)
	interval := l * n / 8
	j := 0

	for i := 0; i < interval*kyb.k; i += interval {
		polyVec[j] = kyb.decode(bytes[i:i+interval], l)
		j++
	}
	return
}

func (kyb *kyber) encodePolyVec(polyVec [][]int, l int) (bytes []byte) {
	for i := 0; i < kyb.k; i++ {
		bytes = append(bytes, kyb.encode(polyVec[i], l)...)
	}
	return
}

func (kyb *kyber) nttPolyVec(polyVec [][]int) {
	for i := 0; i < kyb.k; i++ {
		kyb.ntt(polyVec[i])
	}
}

func (kyb *kyber) invNttPolyVec(polyVec [][]int) {
	for i := 0; i < kyb.k; i++ {
		kyb.invNtt(polyVec[i])
	}
}

func (kyb *kyber) compressPolyVec(polyVec [][]int) {
	for i := 0; i < kyb.k; i++ {

	}
}

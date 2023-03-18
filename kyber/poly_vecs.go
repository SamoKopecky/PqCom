package kyber

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

func modPlusPolyVec(a [][]int) {
	for i := 0; i < K; i++ {
		for j := 0; j < N; j++ {
			a[i][j] = modPlus(a[i][j], Q)
		}
	}
	return
}

func randPolyVec(r []byte, localN *byte, eta int) (vector [][]int) {
	vector = [][]int{}
	for i := 0; i < K; i++ {
		vector = append(vector, randPoly(r, localN, eta))
		*localN += 1
	}
	return
}

func decodePolyVec(bytes []byte, l int) (polyVec [][]int) {
	polyVec = make([][]int, K)
	interval := l * N / 8
	j := 0

	for i := 0; i < interval*K; i += interval {
		polyVec[j] = decode(bytes[i:i+interval], l)
		j++
	}
	return
}

func encodePolyVec(polyVec [][]int, l int) (bytes []byte) {
	for i := 0; i < K; i++ {
		bytes = append(bytes, encode(polyVec[i], l)...)
	}
	return
}

func nttPolyVec(polyVec [][]int) {
	for i := 0; i < K; i++ {
		ntt(polyVec[i])
	}
}

func invNttPolyVec(polyVec [][]int) {
	for i := 0; i < K; i++ {
		invNtt(polyVec[i])
	}
}

func compressPolyVec(polyVec [][]int) {
	for i := 0; i < K; i++ {

	}
}

package dilithium

func polyMul(f, g []int) (h []int) {
	h = make([]int, N)
	for i := 0; i < N; i++ {
		h[i] = (f[i] * g[i]) % Q
	}
	return
}

func polyAdd(f, g []int) (h []int) {
	h = make([]int, N)
	for i := 0; i < N; i++ {
		h[i] = (f[i] + g[i]) % Q
	}
	return
}

func polySub(f, g []int) (h []int) {
	h = make([]int, N)
	for i := 0; i < N; i++ {
		h[i] = (f[i] - g[i]) % Q
	}
	return
}


package dilithium

func (dil *dilithium) polyMul(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] * g[i]) % q
	}
	return
}

func (dil *dilithium) polyAdd(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] + g[i]) % q
	}
	return
}

func (dil *dilithium) polySub(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] - g[i]) % q
	}
	return
}

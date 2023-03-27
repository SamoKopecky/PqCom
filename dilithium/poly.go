package dilithium

func (dil *Dilithium) polyMul(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] * g[i]) % q
	}
	return
}

func (dil *Dilithium) polyAdd(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] + g[i])
	}
	return
}

func (dil *Dilithium) polySub(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] - g[i])
	}
	return
}

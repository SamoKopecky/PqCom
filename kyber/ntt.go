package kyber

func (kyber Kyber) GenerateZetas() (zetas []int) {
	rootOfUnity := 17
	for i := 0; i < 128; i++ {
		zetas = append(zetas, modPow(rootOfUnity, int(bitRev(uint8(i))), kyber.Q))
	}
	return
}

func modPow(base int, power int, modulus int) (result int) {
	// TODO: Make more efficient
	c := 1
	for i := 1; i <= power; i++ {
		c = base * c % modulus
	}
	result = c
	return

}

func bitRev(a uint8) (b uint8) {
	power := 6
	for a != 0 {
		b += uint8((a & 1) << power)
		a = a >> 1
		power -= 1
	}
	return
}

func (kyber Kyber) NTT(poly []int) {
	zetaIndex := 0
	for k := 128; k > 1; k >>= 1 {
		for j := 0; j < 256; j += k * 2 {
			zetaIndex++
			zeta := kyber.Zetas[zetaIndex]
			for i := j; i < j+k; i++ {
				change := zeta * poly[k+i] % kyber.Q
				poly[i+k] = poly[i] - change
				poly[i] += change
			}
		}
	}

	return
}

func (kyber Kyber) InvNTT(poly []int) {
	zetaIndex := 127
	for k := 2; k < 256; k <<= 1 {
		for j := 0; j < 256; j += k * 2 {
			zeta := kyber.Zetas[zetaIndex]
			zetaIndex--
			for i := j; i < j+k; i++ {
				old := poly[i]
				poly[i] += poly[i+k]
				poly[i+k] -= old
				poly[i+k] = zeta * poly[i+k] % kyber.Q
			}
		}
	}

	for i := 0; i < 256; i++ {
		// 3303 == 2^-7 % Q
		poly[i] = (3303 * (int(poly[i]))) % kyber.Q
	}
}

func (kyber Kyber) multiplyPolys(f, g []int, i int, zeta int) (h0, h1 int) {
	h0_1 := (f[i] * g[i]) % kyber.Q
	h0_2 := (f[i+1] * g[i+1]) % kyber.Q
	h0 = (h0_2*zeta + h0_1) % kyber.Q

	h1 = (f[i] * g[i+1]) % kyber.Q
	h1 += (f[i+1] * g[i]) % kyber.Q
	return
}

func (kyber Kyber) polyMulOne(f, g []int) (h []int) {
	h = make([]int, kyber.N)
	zetaIndex := 64
	for i := 0; i < kyber.N; i += 4 {
		zeta := kyber.Zetas[zetaIndex]
		zetaIndex++

		h0, h1 := kyber.multiplyPolys(f, g, i, zeta)
		h[i] = h0
		h[i+1] = h1

		h2, h3 := kyber.multiplyPolys(f, g, i+2, -zeta)
		h[i+2] = h2
		h[i+3] = h3
	}
	return
}

func (kyber Kyber) PointwisePolyMul(f, g [][]int) (h []int) {
	h = make([]int, kyber.N)
	for i := 0; i < kyber.K; i++ {
		h = kyber.PolyAddOne(kyber.polyMulOne(f[i], g[i]), h)
	}
	return
}

func (kyber Kyber) PolyAddOne(f, g []int) (h []int) {
	h = make([]int, kyber.N)
	for i := 0; i < kyber.N; i++ {
		h[i] = (f[i] + g[i]) % kyber.Q
	}
	return
}

func (kyber Kyber) PolySubOne(f, g []int) (h []int) {
	h = make([]int, kyber.N)
	for i := 0; i < kyber.N; i++ {
		h[i] = (f[i] - g[i]) % kyber.Q
	}
	return
}

func (kyber Kyber) PolyAdd(f, g [][]int) (h [][]int) {
	h = make([][]int, kyber.K)
	for i := 0; i < kyber.K; i++ {
		h[i] = kyber.PolyAddOne(f[i], g[i])
	}
	return
}

func (kyber Kyber) ReduceModuloPlusVectors(a [][]int) {
	for i := 0; i < kyber.K; i++ {
		for j := 0; j < kyber.N; j++ {
			a[i][j] = ReduceModuloPlus(a[i][j], kyber.Q)
		}
	}
	return
}

func ReduceModuloPlus(a int, modulo int) (b int) {
	b = a % modulo
	if a < 0 {
		b += modulo
	}
	return
}

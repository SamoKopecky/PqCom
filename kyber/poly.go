package kyber

import "math"

func (kyb *kyber) encode(poly []int, l int) (bytes []byte) {
	bits := []byte{}
	for i := 0; i < 256; i++ {
		for j := 0; j < l; j++ {
			bits = append(bits, byte(poly[i]/int(math.Pow(2, float64(j))))%2)
		}
	}
	var encoded byte
	for i := 0; i < l*256; i += 8 {
		for j := 0; j < 8; j++ {
			encoded += (bits[j+i]) * byte(math.Pow(2, float64(j)))
		}
		bytes = append(bytes, encoded)
		encoded = 0
	}
	return
}

func (kyb *kyber) decode(bytes []byte, l int) (poly []int) {
	bits := kyb.bytesToBits(bytes)
	for i := 0; i < 256; i++ {
		fi := 0
		for j := 0; j < l; j++ {
			fi += int(bits[i*l+j]) * int(math.Pow(2, float64(j)))
		}
		poly = append(poly, fi)
	}
	return
}

func (kyb *kyber) compress(input []int, d int) (compressed []int) {
	for _, v := range input {
		modulo := float64(math.Pow(2, float64(d)))
		parenthesis := modulo / float64(q)
		value := int(math.Round(float64(parenthesis * float64(v))))
		compressed = append(compressed, kyb.modPlus(value, int(modulo)))
	}

	return
}

func (kyb *kyber) decompress(input []int, d int) (decompressed []int) {
	for _, v := range input {
		parenthesis := float64(q) / math.Pow(2, float64(d))
		decompressed = append(decompressed, int(math.Round(parenthesis*float64(v))))
	}
	return
}

func (kyb *kyber) pointWiseMul(f, g []int, i int, zeta int) (h0, h1 int) {
	h0_1 := (f[i] * g[i]) % q
	h0_2 := (f[i+1] * g[i+1]) % q
	h0 = (h0_2*zeta + h0_1) % q

	h1 = (f[i] * g[i+1]) % q
	h1 += (f[i+1] * g[i]) % q
	return
}

func (kyb *kyber) polyMul(f, g []int) (h []int) {
	h = make([]int, n)
	zetaIndex := 64
	for i := 0; i < n; i += 4 {
		zeta := Zetas[zetaIndex]
		zetaIndex++

		h0, h1 := kyb.pointWiseMul(f, g, i, zeta)
		h[i] = h0
		h[i+1] = h1

		h2, h3 := kyb.pointWiseMul(f, g, i+2, -zeta)
		h[i+2] = h2
		h[i+3] = h3
	}
	return
}

func (kyb *kyber) polyAdd(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] + g[i]) % q
	}
	return
}

func (kyb *kyber) polySub(f, g []int) (h []int) {
	h = make([]int, n)
	for i := 0; i < n; i++ {
		h[i] = (f[i] - g[i]) % q
	}
	return
}

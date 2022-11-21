package kyber

import "math"

func encode(poly []int, l int) (bytes []byte) {
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

func decode(bytes []byte, l int) (poly []int) {
	bits := bytesToBits(bytes)
	for i := 0; i < 256; i++ {
		fi := 0
		for j := 0; j < l; j++ {
			fi += int(bits[i*l+j]) * int(math.Pow(2, float64(j)))
		}
		poly = append(poly, fi)
	}
	return
}

func compress(input []int, d int) (compressed []int) {
	for _, v := range input {
		modulo := float64(math.Pow(2, float64(d)))
		parenthesis := modulo / float64(Q)
		value := int(math.Round(float64(parenthesis * float64(v))))
		compressed = append(compressed, modPlus(value, int(modulo)))
	}

	return
}

func decompress(input []int, d int) (decompressed []int) {
	for _, v := range input {
		parenthesis := float64(Q) / math.Pow(2, float64(d))
		decompressed = append(decompressed, int(math.Round(parenthesis*float64(v))))
	}
	return
}

func pointWiseMul(f, g []int, i int, zeta int) (h0, h1 int) {
	h0_1 := (f[i] * g[i]) % Q
	h0_2 := (f[i+1] * g[i+1]) % Q
	h0 = (h0_2*zeta + h0_1) % Q

	h1 = (f[i] * g[i+1]) % Q
	h1 += (f[i+1] * g[i]) % Q
	return
}

func polyMul(f, g []int) (h []int) {
	h = make([]int, N)
	zetaIndex := 64
	for i := 0; i < N; i += 4 {
		zeta := Zetas[zetaIndex]
		zetaIndex++

		h0, h1 := pointWiseMul(f, g, i, zeta)
		h[i] = h0
		h[i+1] = h1

		h2, h3 := pointWiseMul(f, g, i+2, -zeta)
		h[i+2] = h2
		h[i+3] = h3
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

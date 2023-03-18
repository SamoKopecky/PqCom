package kyber

func generateZetas() (zetas []int) {
	rootOfUnity := 17
	for i := 0; i < 128; i++ {
		zetas = append(zetas, modPow(rootOfUnity, int(bitRev(uint8(i))), Q))
	}
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

func modPow(base int, power int, modulus int) (result int) {
	// TODO: Make more efficient
	c := 1
	for i := 1; i <= power; i++ {
		c = base * c % modulus
	}
	result = c
	return
}

var Zetas = []int{
	1, 1729, 2580, 3289, 2642, 630, 1897, 848, 1062, 1919, 193, 797, 2786, 3260, 569, 1746, 296, 2447, 1339, 1476, 3046, 56, 2240, 1333, 1426, 2094, 535, 2882, 2393, 2879, 1974, 821, 289, 331, 3253, 1756, 1197, 2304, 2277, 2055, 650, 1977, 2513, 632, 2865, 33, 1320, 1915, 2319, 1435, 807, 452, 1438, 2868, 1534, 2402, 2647, 2617, 1481, 648, 2474, 3110, 1227, 910, 17, 2761, 583, 2649, 1637, 723, 2288, 1100, 1409, 2662, 3281, 233, 756, 2156, 3015, 3050, 1703, 1651, 2789, 1789, 1847, 952, 1461, 2687, 939, 2308, 2437, 2388, 733, 2337, 268, 641, 1584, 2298, 2037, 3220, 375, 2549, 2090, 1645, 1063, 319, 2773, 757, 2099, 561, 2466, 2594, 2804, 1092, 403, 1026, 1143, 2150, 2775, 886, 1722, 1212, 1874, 1029, 2110, 2935, 885, 2154,
}

func ntt(poly []int) {
	zetaIndex := 0
	for k := 128; k > 1; k >>= 1 {
		for j := 0; j < 256; j += k * 2 {
			zetaIndex++
			zeta := Zetas[zetaIndex]
			for i := j; i < j+k; i++ {
				change := zeta * poly[k+i] % Q
				poly[i+k] = poly[i] - change
				poly[i] += change
			}
		}
	}
	return
}

func invNtt(poly []int) {
	zetaIndex := 127
	for k := 2; k < 256; k <<= 1 {
		for j := 0; j < 256; j += k * 2 {
			zeta := Zetas[zetaIndex]
			zetaIndex--
			for i := j; i < j+k; i++ {
				old := poly[i]
				poly[i] += poly[i+k]
				poly[i+k] -= old
				poly[i+k] = zeta * poly[i+k] % Q
			}
		}
	}

	for i := 0; i < 256; i++ {
		// 3303 == 2^-7 % Q
		poly[i] = (3303 * (int(poly[i]))) % Q
	}
}

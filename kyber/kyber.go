package kyber

import (
	"crypto/rand"
	"math"

	"golang.org/x/crypto/sha3"
)

type Kyber struct {
	Q           int
	N           int
	K           int
	Eta1        int
	Eta2        int
	Du          int
	Dv          int
	Zetas       []int
	MontgomeryR int
}

func (kyber Kyber) Encode(poly []int, l int) (byteStream []byte) {
	bits := []byte{}
	for i := 0; i < 256; i++ {
		for j := 0; j < l; j++ {
			bits = append(bits, byte(poly[i]/int(math.Pow(2, float64(j))))%2)
		}
	}
	encodedNum := byte(0)
	for i := 0; i < l*256; i += 8 {
		for j := 0; j < 8; j++ {
			encodedNum += (bits[j+i]) * byte(math.Pow(2, float64(j)))
		}
		byteStream = append(byteStream, encodedNum)
		encodedNum = 0
	}
	return
}

func (kyber Kyber) Decode(bytes []byte, l int) (poly []int) {
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

func (kyber Kyber) DecodeVectors(bytes []byte, l int) (poly [][]int) {
	poly = make([][]int, kyber.K)
	interval := l * kyber.N / 8
	j := 0

	for i := 0; i < interval*kyber.K; i += interval {
		poly[j] = kyber.Decode(bytes[i:i+interval], l)
		j++
	}
	return
}

func (kyber Kyber) Compress(input []int, d int) (compressed []int) {
	for _, v := range input {
		modulo := float64(math.Pow(2, float64(d)))
		parenthesis := modulo / float64(kyber.Q)
		value := int(math.Round(float64(parenthesis * float64(v))))
		compressed = append(compressed, ReduceModuloPlus(value, int(modulo)))
	}

	return
}

func (kyber Kyber) Decompress(input []int, d int) (decompressed []int) {
	for _, v := range input {
		parenthesis := float64(kyber.Q) / math.Pow(2, float64(d))
		decompressed = append(decompressed, int(math.Round(parenthesis*float64(v))))
	}
	return
}

func CBD(byteStream []byte, eta int) (poly []int) {
	bits := bytesToBits(byteStream)
	for i := 0; i < 256; i++ {
		a := 0
		b := 0
		for j := 0; j < eta; j++ {
			a += int(bits[2*i*eta+j])
		}
		for j := 0; j < eta; j++ {
			b += int(bits[2*i*eta+eta+j])
		}
		poly = append(poly, a-b)
	}
	return
}

func bytesToBits(bytes []byte) (bits []byte) {
	for i := 0; i < len(bytes)*8; i++ {
		bits = append(bits, (bytes[i/8]/byte(math.Pow(2, float64(i%8))))%2)
	}
	return
}

func RandomBytes(size int) (randBytes []byte) {
	randBytes = make([]byte, size)
	rand.Read(randBytes)
	return
}

func (kyber Kyber) Parse(bytes []byte) (ntt_poly []int) {
	j, i := 0, 0
	ntt_poly = make([]int, kyber.N)
	for j < kyber.N {
		reduced := ReduceModuloPlus(int(bytes[i+1]), 16)
		d1 := int(bytes[i]) + 256*reduced
		d2 := int(math.Floor(float64(bytes[i+1]))) + int((16 * bytes[i+2]))
		if d1 < kyber.Q {
			ntt_poly[j] = d1
			j++
		}
		if d2 < kyber.Q && j < kyber.N {
			ntt_poly[j] = d2
			j++
		}
		i = i + 3
	}
	return
}

func PRF(input []byte, N byte, len int) (output []byte) {
	output = make([]byte, len)
	input = append(input, N)
	sha3.ShakeSum256(output, input)
	return
}

func XOF(input []byte, x byte, y byte, len int) (output []byte) {
	output = make([]byte, len)
	input = append(input, x)
	input = append(input, y)
	sha3.ShakeSum256(output, input)
	return
}

func H(input []byte) (hashedBytes [32]byte) {
	return sha3.Sum256(input)
}

func G(input []byte) (first []byte, second []byte) {
	hashedBytes := sha3.Sum512(input)
	first = hashedBytes[:32]
	second = hashedBytes[32:]
	return
}

func KDF(input []byte, output []byte) {
	sha3.ShakeSum256(output, input)
}

func (kyber Kyber) CpapkeKeyGen() (pk []byte, sk []byte) {
	byteStream := RandomBytes(32)
	randomIndex := RandomBytes(1)
	d := make([]byte, 256)
	d = append(d, byteStream[randomIndex[0]%32])
	A_hat := [][][]int{}
	s_hat := [][]int{}
	e_hat := [][]int{}

	rho, sigma := G(d)
	N := byte(0)

	for i := 0; i < kyber.K; i++ {
		A_row := [][]int{}
		for j := 0; j < kyber.K; j++ {
			A_row = append(A_row, kyber.Parse(XOF(rho, byte(j), byte(i), kyber.N*3)))
		}
		A_hat = append(A_hat, A_row)
	}

	s_hat = kyber.randomVectors(&N, sigma, kyber.Eta1)
	e_hat = kyber.randomVectors(&N, sigma, kyber.Eta1)

	for i := 0; i < kyber.K; i++ {
		kyber.NTT(s_hat[i])
		kyber.NTT(e_hat[i])
	}


	t_hat := make([][]int, kyber.K)
	for i := 0; i < kyber.K; i++ {
		t_hat[i] = kyber.PointwisePolyMul(A_hat[i], s_hat)
	}
	t_hat = kyber.PolyAdd(e_hat, t_hat)

	kyber.ReduceModuloPlusVectors(s_hat)
	kyber.ReduceModuloPlusVectors(t_hat)

	sk = make([]byte, 0)
	for i := 0; i < kyber.K; i++ {
		sk = append(sk, kyber.Encode(s_hat[i], 12)...)
	}

	pk = make([]byte, 0)
	for i := 0; i < kyber.K; i++ {
		pk = append(pk, kyber.Encode(t_hat[i], 12)...)
	}
	pk = append(pk, rho...)
	return

}

func (kyber Kyber) CpapkeEnc(pk []byte, m []byte, randomCoins []byte) (c []byte) {
	N := byte(0)
	t_hat := kyber.DecodeVectors(pk, 12)

	A_hat := [][][]int{}
	r_hat := [][]int{}
	e1 := [][]int{}

	rho := pk[len(pk)-32:]

	for i := 0; i < kyber.K; i++ {
		A_row := [][]int{}
		for j := 0; j < kyber.K; j++ {
			A_row = append(A_row, kyber.Parse(XOF(rho, byte(i), byte(j), kyber.N*3)))
		}
		A_hat = append(A_hat, A_row)
	}

	r_hat = kyber.randomVectors(&N, randomCoins, kyber.Eta1)
	e1 = kyber.randomVectors(&N, randomCoins, kyber.Eta2)

	e2 := CBD(PRF(randomCoins, N, 64*kyber.Eta2), kyber.Eta2)

	for i := 0; i < kyber.K; i++ {
		kyber.NTT(r_hat[i])
	}

	u := make([][]int, kyber.K)
	for i := 0; i < kyber.K; i++ {
		u[i] = kyber.PointwisePolyMul(A_hat[i], r_hat)
		kyber.InvNTT(u[i])
	}
	u = kyber.PolyAdd(u, e1)

	parsed_m := kyber.Decompress(kyber.Decode(m, 1), 1)
	v := kyber.PointwisePolyMul(t_hat, r_hat)
	kyber.InvNTT(v)
	v = kyber.PolyAddOne(v, e2)
	v = kyber.PolyAddOne(v, parsed_m)

	// kyber.ReduceModuloPlusVectors(u)
	// for i := 0; i < kyber.N; i++ {
	// 	v[i] = ReduceModuloPlus(v[i], kyber.Q)
	// }

	c1 := []byte{}
	for i := 0; i < kyber.K; i++ {
		c1 = append(c1, kyber.Encode(kyber.Compress(u[i], kyber.Du), kyber.Du)...)
	}

	c2 := kyber.Encode(kyber.Compress(v, kyber.Dv), kyber.Dv)
	c = append(c1, c2...)

	return
}

func (kyber Kyber) randomVectors(N *byte, sigma []byte, eta int) (vector [][]int) {
	vector = [][]int{}
	for i := 0; i < kyber.K; i++ {
		vector = append(vector, CBD(PRF(sigma, *N, 64*eta), eta))
		*N += 1
	}
	return
}

func (kyber Kyber) CpapkeDec(sk []byte, c []byte) (m []byte) {
	c2 := c[kyber.Du*kyber.K*kyber.N/8:]
	u_decoded := kyber.DecodeVectors(c, kyber.Du)

	u_hat := make([][]int, kyber.K)

	for i := 0; i < kyber.K; i++ {
		u_hat[i] = kyber.Decompress(u_decoded[i], kyber.Du)
	}
	for i := 0; i < kyber.K; i++ {
		kyber.NTT(u_hat[i])
	}

	v := kyber.Decompress(kyber.Decode(c2, kyber.Dv), kyber.Dv)
	s_hat := kyber.DecodeVectors(sk, 12)

	s_hat_u_hat := kyber.PointwisePolyMul(s_hat, u_hat)
	kyber.InvNTT(s_hat_u_hat)
	first_m := kyber.PolySubOne(v, s_hat_u_hat)
	
	m = kyber.Encode(kyber.Compress(first_m, 1), 1)
	return
}

// func (kyber Kyber) MontReduce(a int) (b int) {
// 	// TODO: do this better
// 	// 169 == 2^-16 % Q
// 	return a * 169 % kyber.Q
// }

// func (kyber Kyber) MontReduce(a int) (b int) {
// 	// TODO: do this better
// 	mont_mask := 65535
// 	q_inv := 3327
// 	u := ((a & mont_mask) * q_inv) & mont_mask
// 	t := a + u*kyber.Q
// 	t = t >> 16
// 	if t >= kyber.Q {
// 		t = t - kyber.Q
// 	}
// 	return t
// }

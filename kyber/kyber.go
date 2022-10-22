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
	for i := 0; i < l*256; i++ {
		if (i%8 == 0 && i != 0) || i == l*256-1 {
			byteStream = append(byteStream, encodedNum)
			encodedNum = 0
		}
		encodedNum += (bits[i]) * byte(math.Pow(2, float64(i%8)))
	}
	return
}

func (kyber Kyber) Decode(byteStream []byte, l int) (poly []int) {
	bits := bytesToBits(byteStream)
	for i := 0; i < 256; i++ {
		fi := 0
		for j := 0; j < l; j++ {
			fi += int(float64(bits[i*l+j]) * math.Pow(2, float64(j)))
		}
		poly = append(poly, fi)
	}
	return
}

func (kyber Kyber) Compress(input []int, d int) (compressed []int) {
	for _, v := range input {
		modulo := math.Pow(2, float64(d))
		parenthesis := (math.Pow(2, float64(d))) / float64(kyber.Q)
		compressed = append(compressed, int(math.Round(parenthesis*float64(v)))%int(modulo))
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

func CenteredBinomialDistribution(byteStream []byte, eta int) (poly []int) {
	bits := bytesToBits(byteStream)
	for i := 0; i < 256; i++ {
		a := 0
		b := 0
		for j := 0; j < eta-1; j++ {
			a += int(bits[2*i*eta+j])
		}
		for j := 0; j < eta-1; j++ {
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

func (kyber Kyber) Parse(byteStream []byte) (ntt_poly []int) {
	j, i := 0, 0

	for j < kyber.N {
		d1 := int(byteStream[i]) + 256*int(byteStream[i+1]%16)
		d2 := int(math.Floor(float64(byteStream[i+1]))) + int((16 * byteStream[i+2]))
		if d1 < kyber.Q {
			ntt_poly = append(ntt_poly, d1)
			j++
		}
		if d2 < kyber.Q && j < kyber.N {
			ntt_poly = append(ntt_poly, d2)
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
	d := []byte{}
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

	for i := 0; i < kyber.K; i++ {
		s_hat = append(s_hat, CenteredBinomialDistribution(PRF(sigma, N, 64*kyber.Eta1), kyber.Eta1))
		N += 1
	}

	for i := 0; i < kyber.K; i++ {
		e_hat = append(e_hat, CenteredBinomialDistribution(PRF(sigma, N, 64*kyber.Eta1), kyber.Eta1))
		N += 1
	}

	for i := 0; i < kyber.K; i++ {
		kyber.NTT(s_hat[i])
		kyber.NTT(e_hat[i])
	}
	
	t_hat := make([][]int, kyber.K)
	for i := 0; i < kyber.K; i++ {
		t_hat[i] = kyber.PointwisePolyMul(A_hat[i], s_hat)

	}
	t_hat = kyber.PolyAdd(e_hat, t_hat)

	kyber.ReduceModuloPlus(s_hat)
	kyber.ReduceModuloPlus(t_hat)

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

// func (kyber Kyber) MontReduce(a int) (b int) {
// 	// TODO: do this better
// 	// 169 == 2^-16 % Q
// 	return a * 169 % kyber.Q
// }

func (kyber Kyber) MontReduce(a int) (b int) {
	// TODO: do this better
	mont_mask := 65535
	q_inv := 3327
	u := ((a & mont_mask) * q_inv) & mont_mask
	t := a + u*kyber.Q
	t = t >> 16
	if t >= kyber.Q {
		t = t - kyber.Q
	}
	return t
}

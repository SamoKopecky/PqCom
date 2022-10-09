package kyber

import (
	"crypto/rand"
	"math"
)

func Encode(poly []int, l int) (byteStream []byte) {
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

func Decode(byteStream []byte, l int) (poly []int) {
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

func Compress(input []int, d int, q int) (compressed []int) {
	for _, v := range input {
		modulo := math.Pow(2, float64(d))
		parenthesis := (math.Pow(2, float64(d))) / float64(q)
		compressed = append(compressed, int(math.Round(parenthesis*float64(v)))%int(modulo))
	}

	return
}

func Decompress(input []int, d int, q int) (decompressed []int) {
	for _, v := range input {
		parenthesis := float64(q) / math.Pow(2, float64(d))
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

func Parse(byteStream []byte, q int, n int) (ntt_poly []int) {
	j, i := 0, 0

	for j < n {
		d1 := int(byteStream[i]) + 256*int(byteStream[i+1]%16)
		d2 := int(math.Floor(float64(byteStream[i+1]))) + int((16 * byteStream[i+2]))
		if d1 < q {
			ntt_poly = append(ntt_poly, d1)
			j++
		}
		if d2 < q && j < n {
			ntt_poly = append(ntt_poly, d2)
			j++
		}
		i = i + 3
	}
	return
}

func NTT() {}

func InverseNTT() {}

func XOF() {}

func H() {}

func G() {}

func KDF() {}

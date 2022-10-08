package kyber

import (
	"crypto/rand"
	"fmt"
	"math"
)

func converToDec(number string) {

}

func Encode(poly []int, l int) (result []int) {
	var bits []int
	for i := 0; i < 256; i++ {
		for j := 0; j < l; j++ {
			bits = append(bits, (poly[i]/int(math.Pow(2, float64(j))))%2)
		}
	}
	var number int
	for i := 0; i < l*256; i++ {
		if (i%8 == 0 && i != 0) || i == l*256-1 {
			result = append(result, number)
			number = 0
		}
		number += int(bits[i]) * int(math.Pow(2, float64(i%8)))
	}

	return
}

func Decode(byteStream []byte, l int) (result []int) {
	bits := bytesToBits(byteStream)
	for _, v := range bits {
		fmt.Printf("%b", v)
	}
	for i := 0; i < 256; i++ {
		fi := 0
		for j := 0; j < l; j++ {
			fi += int(float64(bits[i*l+j]) * math.Pow(2, float64(j)))
		}
		result = append(result, fi)
	}
	return
}

// TODO: make input an array
func Compress(x int, d int, q int) (result int) {
	modulo := math.Pow(2, float64(d))
	parathesis := (math.Pow(2, float64(d))) / float64(q)
	result = int(math.Round(parathesis*float64(x))) % int(modulo)
	return
}

// TODO: make input an array
func Decompress(x int, d int, q int) (result int) {
	parathesis := float64(q) / math.Pow(2, float64(d))
	result = int(math.Round(parathesis) * float64(x))
	return
}

func CBD(eta int) (result []int) {
	byteStream, _ := RandomBytes(64 * eta)
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
		result = append(result, a-b)
	}
	return
}

func bytesToBits(bytes []byte) (bits []byte) {
	for i := 0; i < len(bytes)*8; i++ {
		bits = append(bits, (bytes[i/8]/byte(math.Pow(2, float64(i%8))))%2)
	}
	return
}

func RandomBytes(size int) (randBytes []byte, err error) {
	randBytes = make([]byte, size)
	_, err = rand.Read(randBytes)
	return
}

func Parse(byteStream []byte, q int, n int) (result []int) {

	result = []int{}
	j, i := 0, 0

	for j < n {
		d1 := int(byteStream[i]) + int(byteStream[i+1]%16)*256
		d2 := int(math.Floor(float64(byteStream[i+1]))) + int((16 * byteStream[i+2]))
		if d1 < q {
			result = append(result, d1)
			j++
		}
		if d2 < q && j < n {
			result = append(result, d2)
			j++
		}
		i = i + 3
	}

	return
}

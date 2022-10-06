package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	byteStream, _ := RandomBytes(512)
	eta := 2
	q := 3323
	n := 256
	l := 10

	resultParsed := Parse(byteStream, q, n)
	print(resultParsed)

	result := CBD(eta)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}

	resultDecode := Decode(l)
	for _, v := range resultDecode {
		fmt.Printf("%d ", v)
	}
}

func Encode()

func Decode(l int) (result []int) {
	byteStream, _ := RandomBytes(32 * l)
	bits := bytesToBits(byteStream)
	for i := 0; i < 256; i++ {
		fi := 0
		for j := 0; j < l-1; j++ {
			fi += int(float64(bits[i*l+j]) * math.Pow(2, float64(j)))
		}
		result = append(result, fi)
	}
	return
}

func Compress(x int, d int, q int) (result int) {
	modulo := math.Pow(2, float64(d))
	parathesis := math.Pow(2, float64(d)) / float64(q)
	result = int(math.Round(parathesis*float64(x))) % int(modulo)
	return
}

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
	for _, v := range bytes {
		for j := 0; j < 8; j++ {
			bits = append(bits, (v>>j)&1)
		}
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

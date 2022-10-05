package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	Parse()
}

func RandomBytes(size int) (randBytes []byte, err error) {
	randBytes = make([]byte, size)
	_, err = rand.Read(randBytes)
	return
}

func Parse() {
	rand.Seed(time.Now().UnixNano())
	myBytes, _ := RandomBytes(512)
	result := []int{}
	n := 256
	j, i := 0, 0
	q := 3329

	for j < n {
		d1 := int(myBytes[i]) + int(myBytes[i+1]%16)*256
		d2 := int(math.Floor(float64(myBytes[i+1]))) + int((16 * myBytes[i+2]))
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

	for i, v := range result {
		fmt.Printf("%dX^%d+", v, i)
	}
	fmt.Printf("\nLen: %d\n", len(result))

}

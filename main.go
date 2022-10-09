package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/SamoKopecky/pqcom/main/kyber"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// l := 10
	byteStream := kyber.RandomBytes(512)
	// eta := 2
	q := 3329
	// d := 10
	n := 256
	// fmt.Printf("%d\n\n", sha3.Sum512(byteStream))

	resultParsed := kyber.Parse(byteStream, q, n)
	for _, v := range resultParsed {
		fmt.Printf("%d ", v)
	}

	// result := kyber.CenteredBinomialDistribution(byteStream, eta)
	// for _, v := range result {
	// 	fmt.Printf("%d ", v)
	// }

	// fmt.Printf("\nLen of stream: %d\n", len(byteStream))
	// for _, v := range byteStream {
	// 	fmt.Printf("%d ", v)
	// }
	// resultDecode := kyber.Decode(byteStream, l)
	// fmt.Printf("\nLen of decode: %d\n", len(resultDecode))
	// for _, v := range resultDecode {
	// 	fmt.Printf("%d ", v)
	// }
	// resultEn`c`ode := kyber.Encode(resultDecode, l)
	// fmt.Printf("\nLen of encode: %d\n", len(resultEncode))
	// for _, v := range resultEncode {
	// 	fmt.Printf("%d ", v)
	// }

	// fmt.Print("Random ints:\n")
	// var randomInts []int
	// for i := 0; i < 256; i++ {
	// 	randomInts = append(randomInts, rand.Intn(q-1))
	// }
	// for _, v := range randomInts {
	// 	fmt.Printf("%d ", v)
	// }

	// fmt.Print("\nCompressed:\n")
	// compressed := kyber.Compress(randomInts, d, q)
	// for _, v := range compressed {
	// 	fmt.Printf("%d ", v)
	// }

	// fmt.Printf("\n\n%f\n\n", math.Round(float64(q)/(math.Pow(2, float64(d+1)))))

	// fmt.Print("\nDecompressed:\n")
	// decompressed := kyber.Decompress(compressed, d, q)
	// for _, v := range decompressed {
	// 	fmt.Printf("%d ", v)
	// }

	// fmt.Print("\nDiff:\n")
	// for i := 0; i < len(compressed); i++ {
	// 	fmt.Printf("%d ", (decompressed[i]-randomInts[i])%q)
	// }

}

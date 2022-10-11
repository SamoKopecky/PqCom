package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// l := 20
	// byteStream := kyber.RandomBytes(32 * l)
	// eta := 2
	// q := 3329
	// d := 10
	// n := 256

	// resultParsed := kyber.Parse(byteStream, q, n)
	// fmt.Printf("%d ", resultParsed)

	// first, second := kyber.G(byteStream)
	// fmt.Printf("\n%d\n", first)
	// fmt.Printf("\n%d\n", second)

	// output := [256]byte{}
	// kyber.KDF(byteStream, output[:])
	// fmt.Printf("%x", output)

	// result := kyber.CenteredBinomialDistribution(byteStream, eta)
	// fmt.Printf("%d ", result)

	// fmt.Printf("\nLen of stream: %d\n", len(byteStream))
	// fmt.Printf("\n%d\n", byteStream)

	// resultDecode := kyber.Decode(byteStream, l)
	// fmt.Printf("\nLen of decode: %d\n", len(resultDecode))
	// fmt.Printf("\n%d\n", resultDecode)

	// resultEncode := kyber.Encode(resultDecode, l)
	// fmt.Printf("\nLen of encode: %d\n", len(resultEncode))
	// fmt.Printf("%d\n", resultEncode)

	// fmt.Print("Random ints:\n")
	// var randomInts []int
	// for i := 0; i < 256; i++ {
	// 	randomInts = append(randomInts, rand.Intn(q-1))
	// }
	// fmt.Printf("%d ", randomInts)

	// fmt.Print("\nCompressed:\n")
	// compressed := kyber.Compress(randomInts, d, q)
	// fmt.Printf("%d ", compressed)

	// fmt.Print("\nDecompressed:\n")
	// decompressed := kyber.Decompress(compressed, d, q)
	// fmt.Printf("%d ", decompressed)

	// fmt.Print("\nDiff:\n")
	// for i := 0; i < len(compressed); i++ {
	// 	fmt.Printf("%d ", (decompressed[i]-randomInts[i])%q)
	// }

}

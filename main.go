package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/SamoKopecky/pqcom/main/kyber"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// l := 3
	// byteStream, _ := kyber.RandomBytes(32*l)
	// eta := 2
	q := 3323
	// n := 256
	// fmt.Printf("%d\n\n", sha3.Sum512(byteStream))

	// resultParsed := kyber.Parse(byteStream, q, n)
	// print(resultParsed)

	// result := kyber.CBD(eta)
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
	// resultEncode := kyber.Encode(resultDecode, l)
	// fmt.Printf("\nLen of encode: %d\n", len(resultEncode))
	// for _, v := range resultEncode {
	// 	fmt.Printf("%d ", v)
	// }
	test := int(1869)
	fmt.Printf("Original: %d\n", test)
	compressed := kyber.Compress(test, 6, q)
	fmt.Printf("Compressed: %d\n", compressed)
	fmt.Printf("Decompressed: %d\n", kyber.Decompress(compressed, 6, q))
}

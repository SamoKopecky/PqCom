package main

import (
	"fmt"
	"math"

	"github.com/SamoKopecky/pqcom/main/kyber"
)

func main() {
	kyber512 := kyber.Kyber{N: 256, K: 2, Q: 3329, Eta1: 3, Eta2: 2, Du: 10, Dv: 4}
	kyber512.MontgomeryR = int(math.Pow(2.0, 16.0)) % kyber512.Q
	kyber512.Zetas = kyber512.GenerateZetas()

	m := kyber.RandomBytes(32)
	r := kyber.RandomBytes(32)

	pk, sk := kyber512.CpapkeKeyGen()
	c := kyber512.CpapkeEnc(pk, m, r)

	
	// fmt.Printf("Cipher text: \n\n%d\n\n", c)

	dec_m := kyber512.CpapkeDec(sk, c)

	fmt.Printf("Message: \n\n%d\n\n", m)
	fmt.Printf("Decrypted message: \n\n%d\n\n", dec_m)



}

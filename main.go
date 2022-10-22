package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/SamoKopecky/pqcom/main/kyber"
)

func main() {
	kyber512 := kyber.Kyber{N: 256, K: 2, Q: 3329, Eta1: 3, Eta2: 2, Du: 10, Dv: 4}
	kyber512.MontgomeryR = int(math.Pow(2.0, 16.0)) % kyber512.Q
	kyber512.Zetas = kyber512.GenerateZetas()
	rand.Seed(time.Now().UnixNano())

	pk, sk := kyber512.CpapkeKeyGen()

	fmt.Printf("%d\n\n", pk)
	fmt.Printf("%d", sk)

}

package main

import (
	"fmt"

	"github.com/SamoKopecky/pqcom/main/kyber"
)

func main() {
	pk, sk := kyber.CcakemKeyGen()
	c, key := kyber.CcakemEnc(pk)
	key2 := kyber.CcakemDec(c, sk)
	fmt.Printf("key       : %d\n", key)
	fmt.Printf("Shared key: %d\n", key2)

	for i := 0; i < 1000; i++ {
		pk, sk = kyber.CcakemKeyGen()
		c, key := kyber.CcakemEnc(pk)
		key2 := kyber.CcakemDec(c, sk)
		if !kyber.BytesEqual(key, key2) {
			fmt.Printf("####### KEYS DONT MATCH #######")
		}
	}

}

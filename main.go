package main

import (
	"fmt"

	"github.com/SamoKopecky/pqcom/main/dilithium"
	"github.com/SamoKopecky/pqcom/main/kyber"
)

func main() {
	print("running 1000 dilithium iterations...\n")
	for i := 0; i < 1000; i++ {
		message := []byte("abc")
		pk, sk := dilithium.KeyGen()
		signature := dilithium.Sign(sk, message)
		verified := dilithium.Verify(pk, message, signature)
		if !verified {
			fmt.Printf("####### SIGNATURES DONT MATCH #######\n")
		}
	}
	print("Done\n")

	pk, sk := kyber.CcakemKeyGen()
	c, key := kyber.CcakemEnc(pk)
	key2 := kyber.CcakemDec(c, sk)
	fmt.Printf("key       : %d\n", key)
	fmt.Printf("Shared key: %d\n", key2)
	print("running 1000 kyber iterations...\n")

	for i := 0; i < 1000; i++ {
		// fmt.Printf("%d\n", i)
		pk, sk = kyber.CcakemKeyGen()
		c, key := kyber.CcakemEnc(pk)
		key2 := kyber.CcakemDec(c, sk)
		if !kyber.BytesEqual(key, key2) {
			fmt.Printf("####### KEYS DONT MATCH #######\n")
		}
	}
	print("Done\n")

}

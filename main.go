package main

import (
	"github.com/SamoKopecky/pqcom/main/dilithium"
)

func main() {
	message := []byte("abc")
	pk, sk := dilithium.KeyGen()
	signature := dilithium.Sign(sk, message)
	verified := dilithium.Verify(pk, message, signature)
	print("\nVerified: ", verified)
	// pk, sk := kyber.CcakemK0eyGen()
	// c, key := kyber.CcakemEnc(pk)
	// key2 := kyber.CcakemDec(c, sk)
	// fmt.Printf("key       : %d\n", key)
	// fmt.Printf("Shared key: %d\n", key2)
	// fmt.Printf("running 1000 iterations")

	// for i := 0; i < 1000; i++ {
	// 	// fmt.Printf("%d\n", i)
	// 	pk, sk = kyber.CcakemKeyGen()
	// 	c, key := kyber.CcakemEnc(pk)
	// 	key2 := kyber.CcakemDec(c, sk)
	// 	if !kyber.BytesEqual(key, key2) {
	// 		fmt.Printf("####### KEYS DONT MATCH #######")
	// 	}
	// }

}

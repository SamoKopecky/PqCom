package main

import "github.com/SamoKopecky/pqcom/main/cmd"

func main() {
	cmd.Execute()
	// dil := dilithium.Dilithium2()
	// pk, sk := dil.KeyGen()
	// s := dil.Sign(sk, []byte("abc"))
	// dil.Verify(pk, []byte("abc"), s)
	// kyb := kyber.Kyber1024()
	// pk, _ := kyb.CcakemKeyGen()
	// kyb.CcakemEnc(pk)
}

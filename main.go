package main

import "github.com/SamoKopecky/pqcom/main/kyber"

func main() {
	// cmd.Execute()
	// dil := dilithium.Dilithium2()
	// _, sk := dil.KeyGen()
	// dil.Sign(sk, []byte("abc"))
	kyb := kyber.Kyber1024()
	pk, _ := kyb.CcakemKeyGen()
	kyb.CcakemEnc(pk)
}

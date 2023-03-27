package main

import "github.com/SamoKopecky/pqcom/main/dilithium"

func main() {
	// cmd.Execute()
	dil := dilithium.Dilithium5()
	_, sk := dil.KeyGen()
	dil.Sign(sk, []byte("abc"))
}

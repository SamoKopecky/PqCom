package main

import "github.com/SamoKopecky/pqcom/main/cmd"

func main() {
	cmd.Execute()
	// kyb := kyber.Kyber512()
	// pk, sk := kyb.CcakemKeyGen()
	// // print("original:\n")
	// // kyb.HashSk(pk)
	// c, k1 := kyb.CcakemEnc(pk)
	// k2 := kyb.CcakemDec(c, sk)
	// // kyb.HashSk(sk)
	// if !kyb.BytesEqual(k1, k2) {
	// 	print("###\n")
	// }
	// c, k1 = kyb.CcakemEnc(pk)
	// k2 = kyb.CcakemDec(c, sk)
	// if !kyb.BytesEqual(k1, k2) {
	// 	print("###\n")
	// }
}

package main

import "github.com/SamoKopecky/pqcom/main/cmd"

func main() {
	cmd.Execute()
	// message := []byte("abc")
	// dil := dilithium.Dilithium5()
	// for i := 0; i < 100; i++ {
	// 	pk, sk := dil.KeyGen()
	// 	signature := dil.Sign(sk, message)
	// 	verified := dil.Verify(pk, message, signature)
	// 	if !verified {
	// 		print("rip")
	// 	}
	// }
	// fmt.Println()
}

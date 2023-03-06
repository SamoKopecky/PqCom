package benchmark

import (
	"fmt"

	"github.com/SamoKopecky/pqcom/main/crypto"
)

type test struct {
	kemAlg crypto.KemAlgorithm
}

func Run(iterations int) {
	fmt.Printf("Running benchmarks for %d iterations...\n", iterations)
	allKems := crypto.GetKemNames()
	allSigns := crypto.GetSignNames()
	for _, name := range allKems {
		kems := crypto.GetKem(name)
		timeKem(useKem, kems.F, name, iterations)
	}
	for _, name := range allSigns {
		signs := crypto.GetSign(name)
		timeSign(useSign, signs.F, name, iterations)
	}
	fmt.Println("Done")
}

func useKem(algs crypto.KemAlgorithm) {
	pk, sk := algs.KeyGen()
	c, _ := algs.Enc(pk)
	algs.Dec(c, sk)
}

func useSign(algs crypto.SignAlgorithm) {
	m := []byte("foo")
	pk, sk := algs.KeyGen()
	sig := algs.Sign(sk, m)
	algs.Verify(pk, m, sig)
}

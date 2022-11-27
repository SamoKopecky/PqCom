package main

import (
	"flag"
	"fmt"

	"github.com/SamoKopecky/pqcom/main/dilithium"
	"github.com/SamoKopecky/pqcom/main/kyber"
	"github.com/cloudflare/circl/sign/dilithium/mode2"
	kyberk2so "github.com/symbolicsoft/kyber-k2so"
)

func main() {
	var iterations int
	flag.IntVar(&iterations, "i", 1000, "set the number of iterations")
	flag.Parse()

	fmt.Printf("Running benchmarks for %d iterations...\n", iterations)
	timeFunction(kyberk2soKyber, iterations)
	timeFunction(myKyber, iterations)
	timeFunction(circlDilithium, iterations)
	timeFunction(myDilithium, iterations)
	print("Done.\n")
}

func circlDilithium() {
	message := []byte("abc")
	signature := make([]byte, mode2.SignatureSize)
	pk, sk, _ := mode2.GenerateKey(nil)
	mode2.SignTo(sk, message, signature)
	mode2.Verify(pk, message, signature)
}

func myDilithium() {
	message := []byte("abc")
	pk, sk := dilithium.KeyGen()
	signature := dilithium.Sign(sk, message)
	dilithium.Verify(pk, message, signature)
}

func kyberk2soKyber() {
	privateKey, publicKey, _ := kyberk2so.KemKeypair512()
	ciphertext, _, _ := kyberk2so.KemEncrypt512(publicKey)
	kyberk2so.KemDecrypt512(ciphertext, privateKey)
}

func myKyber() {
	pk, sk := kyber.CcakemKeyGen()
	c, _ := kyber.CcakemEnc(pk)
	kyber.CcakemDec(c, sk)
}

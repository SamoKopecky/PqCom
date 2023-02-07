package benchmark

import (
	"fmt"

	"github.com/SamoKopecky/pqcom/main/dilithium"
	"github.com/SamoKopecky/pqcom/main/kyber"
	"github.com/cloudflare/circl/kem/kyber/kyber512"
	"github.com/cloudflare/circl/sign/dilithium/mode2"
	kyberk2so "github.com/symbolicsoft/kyber-k2so"
)

func Run(iterations int) {
	fmt.Printf("Running benchmarks for %d iterations...\n", iterations)
	timeFunction(myKyber, iterations)
	timeFunction(myDilithium, iterations)
	timeFunction(kyberk2soKyber, iterations)
	timeFunction(circlKyber, iterations)
	timeFunction(circlDilithium, iterations)
	print("Done.\n")
}

func circlKyber() {
	pk, sk, _ := kyber512.GenerateKeyPair(nil)
	k1 := make([]byte, kyber512.SharedKeySize)
	ct := make([]byte, kyber512.CiphertextSize)
	pk.EncapsulateTo(ct, k1, nil)
	k2 := make([]byte, kyber512.SharedKeySize)
	sk.DecapsulateTo(k2, ct)
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

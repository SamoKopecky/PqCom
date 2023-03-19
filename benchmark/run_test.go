package benchmark

import (
	"testing"

	"github.com/SamoKopecky/pqcom/main/dilithium"
	"github.com/SamoKopecky/pqcom/main/kyber"
)

func TestMyKyber(t *testing.T) {
	kyb := kyber.Kyber512()
	pk, sk := kyb.CcakemKeyGen()
	c, k1 := kyb.CcakemEnc(pk)
	k2 := kyb.CcakemDec(c, sk)
	if !kyb.BytesEqual(k1, k2) {
		t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
	}
}

func TestMyDilithium(t *testing.T) {
	message := []byte("abc")
	dil := dilithium.Dilithium2()
	pk, sk := dil.KeyGen()
	signature := dil.Sign(sk, message)
	verified := dil.Verify(pk, message, signature)
	if !verified {
		t.Fatalf("signature needs to be verified\nverified: %t", verified)
	}
}

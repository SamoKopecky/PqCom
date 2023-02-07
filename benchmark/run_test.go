package benchmark

import (
	"testing"

	"github.com/SamoKopecky/pqcom/main/dilithium"
	"github.com/SamoKopecky/pqcom/main/kyber"
)

func TestMyKyber(t *testing.T) {
	pk, sk := kyber.CcakemKeyGen()
	c, k1 := kyber.CcakemEnc(pk)
	k2 := kyber.CcakemDec(c, sk)
	if !kyber.BytesEqual(k1, k2) {
		t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
	}
}

func TestMyDilithium(t *testing.T) {
	message := []byte("abc")
	pk, sk := dilithium.KeyGen()
	signature := dilithium.Sign(sk, message)
	verified := dilithium.Verify(pk, message, signature)
	if !verified {
		t.Fatalf("signature needs to be verified\nverified: %t", verified)
	}
}

package benchmark

import (
	"testing"

	"github.com/SamoKopecky/pqcom/main/dilithium"
	"github.com/SamoKopecky/pqcom/main/kyber"
)

func TestPqComKyber512(t *testing.T) {
	testKyber(kyber.Kyber512(), t)
}
func TestPqComKyber768(t *testing.T) {
	testKyber(kyber.Kyber768(), t)
}
func TestPqComKyber1024(t *testing.T) {
	testKyber(kyber.Kyber1024(), t)
}

func testKyber(kyb kyber.Kyber, t *testing.T) {
	pk, sk := kyb.CcakemKeyGen()
	for i := 0; i < 100; i++ {
		c, k1 := kyb.CcakemEnc(pk)
		k2 := kyb.CcakemDec(c, sk)
		if !kyb.BytesEqual(k1, k2) {
			t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
		}
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

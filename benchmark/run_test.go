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

func TestPqComKyber512SameKeys(t *testing.T) {
	testKyberSameKeys(kyber.Kyber512(), t)
}
func TestPqComKyber768SameKeys(t *testing.T) {
	testKyberSameKeys(kyber.Kyber768(), t)
}
func TestPqComKyber1024SameKeys(t *testing.T) {
	testKyberSameKeys(kyber.Kyber1024(), t)
}

func testKyber(kyb kyber.Kyber, t *testing.T) {
	for i := 0; i < 50; i++ {
		pk, sk := kyb.CcakemKeyGen()
		c, k1 := kyb.CcakemEnc(pk)
		k2 := kyb.CcakemDec(c, sk)
		if !kyb.BytesEqual(k1, k2) {
			t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
		}
	}
}

func testKyberSameKeys(kyb kyber.Kyber, t *testing.T) {
	pk, sk := kyb.CcakemKeyGen()
	for i := 0; i < 50; i++ {
		c, k1 := kyb.CcakemEnc(pk)
		k2 := kyb.CcakemDec(c, sk)
		if !kyb.BytesEqual(k1, k2) {
			t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
		}
	}
}

func TestPqComDilithium2(t *testing.T) {
	testDilithum(dilithium.Dilithium2(), t)
}
func TestPqComDilithium3(t *testing.T) {
	testDilithum(dilithium.Dilithium3(), t)
}
func TestPqComDilithium5(t *testing.T) {
	testDilithum(dilithium.Dilithium5(), t)
}

func TestPqComDilithium2SameKeys(t *testing.T) {
	testDilithumSameKeys(dilithium.Dilithium2(), t)
}
func TestPqComDilithium3SameKeys(t *testing.T) {
	testDilithumSameKeys(dilithium.Dilithium3(), t)
}
func TestPqComDilithium5SameKeys(t *testing.T) {
	testDilithumSameKeys(dilithium.Dilithium5(), t)
}

func testDilithumSameKeys(dil dilithium.Dilithium, t *testing.T) {
	message := []byte("foo")
	pk, sk := dil.KeyGen()
	for i := 0; i < 50; i++ {
		signature := dil.Sign(sk, message)
		verified := dil.Verify(pk, message, signature)
		if !verified {
			t.Fatalf("signature needs to be verified\nverified: %t", verified)
		}
	}
}

func testDilithum(dil dilithium.Dilithium, t *testing.T) {
	message := []byte("bar")
	for i := 0; i < 50; i++ {
		pk, sk := dil.KeyGen()
		signature := dil.Sign(sk, message)
		verified := dil.Verify(pk, message, signature)
		if !verified {
			t.Fatalf("signature needs to be verified\nverified: %t", verified)
		}
	}
}

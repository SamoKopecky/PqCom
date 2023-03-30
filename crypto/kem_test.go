package crypto_test

import (
	"strings"
	"testing"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/SamoKopecky/pqcom/main/kyber"
)

const testIterations = 50

func BenchmarkKem(b *testing.B) {
	for k, v := range crypto.Kems {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runKem(v, b)
			}
		})
	}
}

func BenchmarkPqComKem(b *testing.B) {
	for k, v := range crypto.Kems {
		if !strings.Contains(k, "PqCom") {
			continue
		}
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runKem(v, b)
			}
		})
	}
}

func BenchmarkPqComKem1024(b *testing.B) {
	for k, v := range crypto.Kems {
		if !strings.Contains(k, "PqComKyber1024") {
			continue
		}
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runKem(v, b)
			}
		})
	}
}

func BenchmarkKemKeyGen(b *testing.B) {
	for k, v := range crypto.Kems {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v.KeyGen()
			}
		})
	}
}

func BenchmarkKemEnc(b *testing.B) {
	var pk []byte
	for k, v := range crypto.Kems {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				pk, _ = v.KeyGen()
				b.StartTimer()
				v.Enc(pk)
			}
		})
	}
}

func BenchmarkKemDec(b *testing.B) {
	var pk, sk, c []byte
	for k, v := range crypto.Kems {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				pk, sk = v.KeyGen()
				c, _ = v.Enc(pk)
				b.StartTimer()
				v.Dec(c, sk)
			}
		})
	}
}

func runKem(alg crypto.KemAlgorithm, b *testing.B) {
	pk, sk := alg.KeyGen()
	c, k1 := alg.Enc(pk)
	k2 := alg.Dec(c, sk)
	if !kyber.BytesEqual(k1, k2) {
		b.Fatalf("keys need to equal\n%d\n%d", k1, k2)
	}
}

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
	for i := 0; i < testIterations; i++ {
		pk, sk := kyb.CcakemKeyGen()
		c, k1 := kyb.CcakemEnc(pk)
		k2 := kyb.CcakemDec(c, sk)
		if !kyber.BytesEqual(k1, k2) {
			t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
		}
	}
}

func testKyberSameKeys(kyb kyber.Kyber, t *testing.T) {
	pk, sk := kyb.CcakemKeyGen()
	for i := 0; i < testIterations; i++ {
		c, k1 := kyb.CcakemEnc(pk)
		k2 := kyb.CcakemDec(c, sk)
		if !kyber.BytesEqual(k1, k2) {
			t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
		}
	}
}

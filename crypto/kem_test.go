package crypto_test

import (
	"strings"
	"testing"

	"github.com/SamoKopecky/pqcom/main/common"
	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/SamoKopecky/pqcom/main/kyber"
)

const kyberIterations = 50

func BenchmarkKem(b *testing.B) {
	for k, v := range crypto.Kems {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runKem(v, b)
			}
		})
	}
}

func BenchmarkPqComKyberAll(b *testing.B) {
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

func BenchmarkPqComKyber1024(b *testing.B) {
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
	c, _ := alg.Enc(pk)
	_ = alg.Dec(c, sk)
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
	for i := 0; i < kyberIterations; i++ {
		pk, sk := kyb.CcakemKeyGen()
		c, k1 := kyb.CcakemEnc(pk)
		k2 := kyb.CcakemDec(c, sk)
		if !common.BytesEqual(k1, k2) {
			t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
		}
	}
}

func testKyberSameKeys(kyb kyber.Kyber, t *testing.T) {
	pk, sk := kyb.CcakemKeyGen()
	for i := 0; i < kyberIterations; i++ {
		c, k1 := kyb.CcakemEnc(pk)
		k2 := kyb.CcakemDec(c, sk)
		if !common.BytesEqual(k1, k2) {
			t.Fatalf("keys need to equal\n%d\n%d", k1, k2)
		}
	}
}

package crypto_test

import (
	"strings"
	"testing"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/SamoKopecky/pqcom/main/dilithium"
)

const dilIterations = 50

func BenchmarkSignature(b *testing.B) {
	for k, v := range crypto.Signatures {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runSign(v, b)
			}
		})
	}
}

func BenchmarkPqComDilAll(b *testing.B) {
	for k, v := range crypto.Signatures {
		if !strings.Contains(k, "PqCom") {
			continue
		}
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runSign(v, b)
			}
		})
	}
}

func BenchmarkPqComDil5(b *testing.B) {
	for k, v := range crypto.Signatures {
		if !strings.Contains(k, "PqComDilithium5") {
			continue
		}
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runSign(v, b)
			}
		})
	}
}

func BenchmarkPqComDil2(b *testing.B) {
	for k, v := range crypto.Signatures {
		if !strings.Contains(k, "PqComDilithium2") {
			continue
		}
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runSign(v, b)
			}
		})
	}
}

func BenchmarkSignKeyGen(b *testing.B) {
	for k, v := range crypto.Signatures {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v.KeyGen()
			}
		})
	}
}

func BenchmarkSignSign(b *testing.B) {
	var sk []byte
	for k, v := range crypto.Signatures {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				_, sk = v.KeyGen()
				b.StartTimer()
				v.Sign(sk, []byte("foo"))
			}
		})
	}
}

func BenchmarkSignVerify(b *testing.B) {
	var pk, sk, s []byte
	for k, v := range crypto.Signatures {
		b.Run(k, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				pk, sk = v.KeyGen()
				s = v.Sign(sk, []byte("foo"))
				b.StartTimer()
				v.Verify(pk, []byte("foo"), s)
			}
		})
	}
}

func runSign(alg crypto.SignAlgorithm, b *testing.B) {
	pk, sk := alg.KeyGen()
	s := alg.Sign(sk, []byte("foo"))
	verified := alg.Verify(pk, []byte("foo"), s)
	b.StopTimer()
	if !verified {
		b.Fatalf("signature needs to be verified\nverified: %t", verified)
	}
	b.StartTimer()
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
	for i := 0; i < dilIterations; i++ {
		signature := dil.Sign(sk, message)
		verified := dil.Verify(pk, message, signature)
		if !verified {
			t.Fatalf("signature needs to be verified\nverified: %t", verified)
		}
	}
}

func testDilithum(dil dilithium.Dilithium, t *testing.T) {
	message := []byte("bar")
	for i := 0; i < dilIterations; i++ {
		pk, sk := dil.KeyGen()
		signature := dil.Sign(sk, message)
		verified := dil.Verify(pk, message, signature)
		if !verified {
			t.Fatalf("signature needs to be verified\nverified: %t", verified)
		}
	}
}

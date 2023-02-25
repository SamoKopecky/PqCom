package sign

import "github.com/SamoKopecky/pqcom/main/dilithium"

type MyDilithium struct{}

func (MyDilithium) KeyGen() (pk, sk []byte) {
	return dilithium.KeyGen()
}

func (MyDilithium) Sign(sk, msg []byte) (signature []byte) {
	return dilithium.Sign(sk, msg)
}

func (MyDilithium) Verify(pk, msg, signature []byte) bool {
	return dilithium.Verify(pk, msg, signature)
}

func (MyDilithium) SignLen() (signLen int) {
	return 2420
}

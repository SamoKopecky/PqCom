package sign

import "github.com/SamoKopecky/pqcom/main/dilithium"

type PqComDilithium2 struct{}

func (PqComDilithium2) KeyGen() (pk, sk []byte) {
	dil := dilithium.Dilithium2()
	return dil.KeyGen()
}

func (PqComDilithium2) Sign(sk, msg []byte) (signature []byte) {
	dil := dilithium.Dilithium2()
	return dil.Sign(sk, msg)
}

func (PqComDilithium2) Verify(pk, msg, signature []byte) bool {
	dil := dilithium.Dilithium2()
	return dil.Verify(pk, msg, signature)
}

func (PqComDilithium2) SignLen() (signLen int) {
	// TODO: Put as constant or something
	return 2420
}

func (PqComDilithium2) PkLen() (signLen int) {
	return 1312
}

func (PqComDilithium2) SkLen() (signLen int) {
	return 2528
}

func (PqComDilithium2) Id() uint8 {
	return 0
}

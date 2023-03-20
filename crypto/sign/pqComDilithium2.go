package sign

import "github.com/SamoKopecky/pqcom/main/dilithium"

type PqComDilithium2 struct{}

var dil2 = dilithium.Dilithium2()

func (PqComDilithium2) KeyGen() (pk, sk []byte) {
	return dil2.KeyGen()
}

func (PqComDilithium2) Sign(sk, msg []byte) (signature []byte) {
	return dil2.Sign(sk, msg)
}

func (PqComDilithium2) Verify(pk, msg, signature []byte) bool {
	return dil2.Verify(pk, msg, signature)
}

func (PqComDilithium2) SignLen() (signLen int) {
	return dil2.SigSize
}

func (PqComDilithium2) PkLen() (signLen int) {
	return dil2.PkSize
}

func (PqComDilithium2) SkLen() (signLen int) {
	return dil2.SkSize
}

func (PqComDilithium2) Id() uint8 {
	return 0
}

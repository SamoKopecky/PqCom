package sign

import "github.com/SamoKopecky/pqcom/main/dilithium"

type PqComDilithium3 struct{}

var dil3 = dilithium.Dilithium3()

func (PqComDilithium3) KeyGen() (pk, sk []byte) {
	return dil3.KeyGen()
}

func (PqComDilithium3) Sign(sk, msg []byte) (signature []byte) {
	return dil3.Sign(sk, msg)
}

func (PqComDilithium3) Verify(pk, msg, signature []byte) bool {
	return dil3.Verify(pk, msg, signature)
}

func (PqComDilithium3) SignLen() (signLen int) {
	return dil3.SigSize
}

func (PqComDilithium3) PkLen() (signLen int) {
	return dil3.PkSize
}

func (PqComDilithium3) SkLen() (signLen int) {
	return dil3.SkSize
}

func (PqComDilithium3) Id() uint8 {
	return 4
}

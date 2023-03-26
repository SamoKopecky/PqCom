package sign

import "github.com/SamoKopecky/pqcom/main/dilithium"

type PqComDilithium5 struct{}

var dil5 = dilithium.Dilithium5()

func (PqComDilithium5) KeyGen() (pk, sk []byte) {
	return dil5.KeyGen()
}

func (PqComDilithium5) Sign(sk, msg []byte) (signature []byte) {
	return dil5.Sign(sk, msg)
}

func (PqComDilithium5) Verify(pk, msg, signature []byte) bool {
	return dil5.Verify(pk, msg, signature)
}

func (PqComDilithium5) SignLen() (signLen int) {
	return dil5.SigSize
}

func (PqComDilithium5) PuKLen() (signLen int) {
	return dil5.PkSize
}

func (PqComDilithium5) PrKLen() (signLen int) {
	return dil5.SkSize
}

func (PqComDilithium5) Id() uint8 {
	return 5
}

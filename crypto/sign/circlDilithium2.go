package sign

import (
	"github.com/cloudflare/circl/sign/dilithium/mode2"
)

type CirclDilithium2 struct{}

func (CirclDilithium2) KeyGen() ([]byte, []byte) {
	pkBytes := make([]byte, mode2.PublicKeySize)
	skBytes := make([]byte, mode2.PrivateKeySize)
	pk, sk, _ := mode2.GenerateKey(nil)
	pk.Pack((*[mode2.PublicKeySize]byte)(pkBytes))
	sk.Pack((*[mode2.PrivateKeySize]byte)(skBytes))
	return pkBytes, skBytes
}

func (CirclDilithium2) Sign(skBytes, msg []byte) []byte {
	signature := make([]byte, mode2.SignatureSize)
	sk := &mode2.PrivateKey{}
	sk.Unpack((*[mode2.PrivateKeySize]byte)(skBytes))
	mode2.SignTo(sk, msg, signature)
	return signature
}

func (CirclDilithium2) Verify(pkBytes, msg, signature []byte) bool {
	pk := &mode2.PublicKey{}
	pk.Unpack((*[mode2.PublicKeySize]byte)(pkBytes))
	return mode2.Verify(pk, msg, signature)
}

func (CirclDilithium2) SignLen() (signLen int) {
	return mode2.SignatureSize
}

func (CirclDilithium2) PkLen() (signLen int) {
	return mode2.PublicKeySize
}

func (CirclDilithium2) SkLen() (signLen int) {
	return mode2.PrivateKeySize
}

func (CirclDilithium2) Id() uint8 {
	return 1
}

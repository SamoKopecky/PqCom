package sign

import (
	"github.com/cloudflare/circl/sign/dilithium/mode2"
)

type CirclDilithium struct{}

func (CirclDilithium) KeyGen() ([]byte, []byte) {
	pkBytes := make([]byte, mode2.PublicKeySize)
	skBytes := make([]byte, mode2.PrivateKeySize)
	pk, sk, _ := mode2.GenerateKey(nil)
	pk.Pack((*[1312]byte)(pkBytes))
	sk.Pack((*[2528]byte)(skBytes))
	return pkBytes, skBytes
}

func (CirclDilithium) Sign(skBytes, msg []byte) ([]byte) {
	signature := make([]byte, mode2.SignatureSize)
	sk := &mode2.PrivateKey{}
	sk.Unpack((*[2528]byte)(skBytes))
	mode2.SignTo(sk, msg, signature)
	return signature
}

func (CirclDilithium) Verify(pkBytes, msg, signature []byte) bool {
	pk := &mode2.PublicKey{}
	pk.Unpack((*[1312]byte)(pkBytes))
	return mode2.Verify(pk, msg, signature)
}

func (CirclDilithium) SignLen() (signLen int) {
	return mode2.SignatureSize
}

func (CirclDilithium) PkLen() (signLen int) {
	return mode2.PublicKeySize
}

func (CirclDilithium) SkLen() (signLen int) {
	return mode2.PrivateKeySize
}
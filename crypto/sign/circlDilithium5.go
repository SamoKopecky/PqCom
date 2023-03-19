package sign

import (
	"github.com/cloudflare/circl/sign/dilithium/mode5"
)

type CirclDilithium5 struct{}

func (CirclDilithium5) KeyGen() ([]byte, []byte) {
	pkBytes := make([]byte, mode5.PublicKeySize)
	skBytes := make([]byte, mode5.PrivateKeySize)
	pk, sk, _ := mode5.GenerateKey(nil)
	pk.Pack((*[mode5.PublicKeySize]byte)(pkBytes))
	sk.Pack((*[mode5.PrivateKeySize]byte)(skBytes))
	return pkBytes, skBytes
}

func (CirclDilithium5) Sign(skBytes, msg []byte) []byte {
	signature := make([]byte, mode5.SignatureSize)
	sk := &mode5.PrivateKey{}
	sk.Unpack((*[mode5.PrivateKeySize]byte)(skBytes))
	mode5.SignTo(sk, msg, signature)
	return signature
}

func (CirclDilithium5) Verify(pkBytes, msg, signature []byte) bool {
	pk := &mode5.PublicKey{}
	pk.Unpack((*[mode5.PublicKeySize]byte)(pkBytes))
	return mode5.Verify(pk, msg, signature)
}

func (CirclDilithium5) SignLen() (signLen int) {
	return mode5.SignatureSize
}

func (CirclDilithium5) PkLen() (signLen int) {
	return mode5.PublicKeySize
}

func (CirclDilithium5) SkLen() (signLen int) {
	return mode5.PrivateKeySize
}

func (CirclDilithium5) Id() uint8 {
	return 3
}

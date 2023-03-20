package sign

import "github.com/cloudflare/circl/sign/dilithium/mode3"

type CirclDilithium3 struct{}

func (CirclDilithium3) KeyGen() ([]byte, []byte) {
	pkBytes := make([]byte, mode3.PublicKeySize)
	skBytes := make([]byte, mode3.PrivateKeySize)
	pk, sk, _ := mode3.GenerateKey(nil)
	pk.Pack((*[mode3.PublicKeySize]byte)(pkBytes))
	sk.Pack((*[mode3.PrivateKeySize]byte)(skBytes))
	return pkBytes, skBytes
}

func (CirclDilithium3) Sign(skBytes, msg []byte) []byte {
	signature := make([]byte, mode3.SignatureSize)
	sk := &mode3.PrivateKey{}
	sk.Unpack((*[mode3.PrivateKeySize]byte)(skBytes))
	mode3.SignTo(sk, msg, signature)
	return signature
}

func (CirclDilithium3) Verify(pkBytes, msg, signature []byte) bool {
	pk := &mode3.PublicKey{}
	pk.Unpack((*[mode3.PublicKeySize]byte)(pkBytes))
	return mode3.Verify(pk, msg, signature)
}

func (CirclDilithium3) SignLen() (signLen int) {
	return mode3.SignatureSize
}

func (CirclDilithium3) PkLen() (signLen int) {
	return mode3.PublicKeySize
}

func (CirclDilithium3) SkLen() (signLen int) {
	return mode3.PrivateKeySize
}

func (CirclDilithium3) Id() uint8 {
	return 2
}

package kem

import "github.com/cloudflare/circl/kem/kyber/kyber1024"

type CirclKyber1024 struct{}

func (CirclKyber1024) KeyGen() ([]byte, []byte) {
	pkBytes := make([]byte, kyber1024.PublicKeySize)
	skBytes := make([]byte, kyber1024.PrivateKeySize)
	pk, sk, _ := kyber1024.GenerateKeyPair(nil)
	pk.Pack(pkBytes)
	sk.Pack(skBytes)
	return pkBytes, skBytes
}

func (CirclKyber1024) Enc(pkBytes []byte) (c, key []byte) {
	k1 := make([]byte, kyber1024.SharedKeySize)
	ct := make([]byte, kyber1024.CiphertextSize)
	pk := kyber1024.PublicKey{}
	pk.Unpack(pkBytes)
	pk.EncapsulateTo(ct, k1, nil)
	return ct, k1
}

func (CirclKyber1024) Dec(c, skBytes []byte) (key []byte) {
	sk := kyber1024.PrivateKey{}
	k2 := make([]byte, kyber1024.SharedKeySize)
	sk.Unpack(skBytes)
	sk.DecapsulateTo(k2, c)
	return k2
}

func (CirclKyber1024) EkLen() int {
	return kyber1024.PublicKeySize
}

func (CirclKyber1024) Id() uint8 {
	return 3
}

func (CirclKyber1024) CLen() int {
	print(kyber1024.Scheme().CiphertextSize())
	return kyber1024.Scheme().CiphertextSize()
}

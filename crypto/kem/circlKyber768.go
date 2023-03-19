package kem

import (
	"github.com/cloudflare/circl/kem/kyber/kyber768"
)

type CirclKyber768 struct{}

func (CirclKyber768) KeyGen() ([]byte, []byte) {
	pkBytes := make([]byte, kyber768.PublicKeySize)
	skBytes := make([]byte, kyber768.PrivateKeySize)
	pk, sk, _ := kyber768.GenerateKeyPair(nil)
	pk.Pack(pkBytes)
	sk.Pack(skBytes)
	return pkBytes, skBytes
}

func (CirclKyber768) Enc(pkBytes []byte) (c, key []byte) {
	k1 := make([]byte, kyber768.SharedKeySize)
	ct := make([]byte, kyber768.CiphertextSize)
	pk := kyber768.PublicKey{}
	pk.Unpack(pkBytes)
	pk.EncapsulateTo(ct, k1, nil)
	return ct, k1
}

func (CirclKyber768) Dec(c, skBytes []byte) (key []byte) {
	sk := kyber768.PrivateKey{}
	k2 := make([]byte, kyber768.SharedKeySize)
	sk.Unpack(skBytes)
	sk.DecapsulateTo(k2, c)
	return k2
}

func (CirclKyber768) EkLen() int {
	return kyber768.PublicKeySize
}

func (CirclKyber768) Id() uint8 {
	return 2
}

func (CirclKyber768) CLen() int {
	print(kyber768.Scheme().CiphertextSize())
	return kyber768.Scheme().CiphertextSize()
}

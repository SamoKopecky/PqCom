package kem

import (
	"github.com/cloudflare/circl/kem/kyber/kyber512"
)

type CirclKyber struct{}

func (CirclKyber) KeyGen() ([]byte, []byte) {
	pkBytes := make([]byte, kyber512.PublicKeySize)
	skBytes := make([]byte, kyber512.PrivateKeySize)
	pk, sk, _ := kyber512.GenerateKeyPair(nil)
	pk.Pack(pkBytes)
	sk.Pack(skBytes)
	return pkBytes, skBytes
}

func (CirclKyber) Enc(pkBytes []byte) (c, key []byte) {
	k1 := make([]byte, kyber512.SharedKeySize)
	ct := make([]byte, kyber512.CiphertextSize)
	pk := kyber512.PublicKey{}
	pk.Unpack(pkBytes)
	pk.EncapsulateTo(ct, k1, nil)
	return ct, k1
}

func (CirclKyber) Dec(c, skBytes []byte) (key []byte) {
	sk := kyber512.PrivateKey{}
	k2 := make([]byte, kyber512.SharedKeySize)
	sk.Unpack(skBytes)
	sk.DecapsulateTo(k2, c)
	return k2
}

func (CirclKyber) EkLen() int {
	return kyber512.PublicKeySize
}

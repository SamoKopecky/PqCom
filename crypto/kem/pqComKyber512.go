package kem

import (
	"github.com/SamoKopecky/pqcom/main/kyber"
	"github.com/cloudflare/circl/kem/kyber/kyber512"
)

type PqComKyber512 struct{}

func (PqComKyber512) KeyGen() (pk, sk []byte) {
	kyb := kyber.Kyber512()
	return kyb.CcakemKeyGen()
}

func (PqComKyber512) Enc(pk []byte) (c, key []byte) {
	kyb := kyber.Kyber512()
	return kyb.CcakemEnc(pk)
}

func (PqComKyber512) Dec(c, sk []byte) (key []byte) {
	kyb := kyber.Kyber512()
	return kyb.CcakemDec(c, sk)
}

func (PqComKyber512) EkLen() int {
	// TODO: Put as constant or something
	return 800
}

func (PqComKyber512) Id() uint8 {
	return 0
}

func (PqComKyber512) CLen() int {
	return kyber512.Scheme().CiphertextSize()
}

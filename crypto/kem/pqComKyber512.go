package kem

import (
	"github.com/SamoKopecky/pqcom/main/kyber"
	"github.com/cloudflare/circl/kem/kyber/kyber512"
)

type PqComKyber512 struct{}

func (PqComKyber512) KeyGen() (pk, sk []byte) {
	return kyber.CcakemKeyGen()
}

func (PqComKyber512) Enc(pk []byte) (c, key []byte) {
	return kyber.CcakemEnc(pk)
}

func (PqComKyber512) Dec(c, sk []byte) (key []byte) {
	return kyber.CcakemDec(c, sk)
}

func (PqComKyber512) EkLen() int {
	return 800
}

func (PqComKyber512) Id() uint8 {
	return 0
}

func (PqComKyber512) CLen() int {
	return kyber512.Scheme().CiphertextSize()
}

package kem

import "github.com/SamoKopecky/pqcom/main/kyber"

type PqComKyber512 struct{}

var kyb512 = kyber.Kyber512()

func (PqComKyber512) KeyGen() (pk, sk []byte) {
	return kyb512.CcakemKeyGen()
}

func (PqComKyber512) Enc(pk []byte) (c, key []byte) {
	return kyb512.CcakemEnc(pk)
}

func (PqComKyber512) Dec(c, sk []byte) (key []byte) {
	return kyb512.CcakemDec(c, sk)
}

func (PqComKyber512) EkLen() int {
	return kyb512.PkSize
}

func (PqComKyber512) Id() uint8 {
	return 0
}

func (PqComKyber512) CLen() int {
	return kyb512.CtSize
}

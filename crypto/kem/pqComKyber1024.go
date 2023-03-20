package kem

import "github.com/SamoKopecky/pqcom/main/kyber"

var kyb1024 = kyber.Kyber1024()

type PqComKyber1024 struct{}

func (PqComKyber1024) KeyGen() (pk, sk []byte) {
	return kyb1024.CcakemKeyGen()
}

func (PqComKyber1024) Enc(pk []byte) (c, key []byte) {
	return kyb1024.CcakemEnc(pk)
}

func (PqComKyber1024) Dec(c, sk []byte) (key []byte) {
	return kyb1024.CcakemDec(c, sk)
}

func (PqComKyber1024) EkLen() int {
	return kyb1024.PkSize
}

func (PqComKyber1024) Id() uint8 {
	return 5
}

func (PqComKyber1024) CLen() int {
	return kyb1024.CtSize
}

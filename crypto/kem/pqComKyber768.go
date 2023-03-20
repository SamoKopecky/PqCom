package kem

import "github.com/SamoKopecky/pqcom/main/kyber"

type PqComKyber768 struct{}

var kyb768 = kyber.Kyber768()

func (PqComKyber768) KeyGen() (pk, sk []byte) {
	return kyb768.CcakemKeyGen()
}

func (PqComKyber768) Enc(pk []byte) (c, key []byte) {
	return kyb768.CcakemEnc(pk)
}

func (PqComKyber768) Dec(c, sk []byte) (key []byte) {
	return kyb768.CcakemDec(c, sk)
}

func (PqComKyber768) EkLen() int {
	return kyb768.PkSize
}

func (PqComKyber768) Id() uint8 {
	return 4
}

func (PqComKyber768) CLen() int {
	return kyb768.CtSize
}

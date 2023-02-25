package kem

import "github.com/SamoKopecky/pqcom/main/kyber"

type MyKyber struct{}

func (MyKyber) KeyGen() (pk, sk []byte) {
	return kyber.CcakemKeyGen()
}

func (MyKyber) Enc(pk []byte) (c, key []byte) {
	return kyber.CcakemEnc(pk)
}

func (MyKyber) Dec(c, sk []byte) (key []byte) {
	return kyber.CcakemDec(c, sk)
}

func (MyKyber) EkLen() int {
	return 800
}

package crypto

import "github.com/SamoKopecky/pqcom/main/crypto/kem"

var kems = map[string]KemAlgorithm{
	"PqComKyber512": &kem.PqComKyber512{},
	"CirclKyber512": &kem.CirclKyber512{},
}

type KemAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Dec(c, sk []byte) (key []byte)
	Enc(pk []byte) (c, key []byte)
	EkLen() (ekLen int)
}

type Kem struct {
	Id        int
	Functions KemAlgorithm
}

func GetKem(kemName string) Kem {
	id := getRowIndex(kems, kemName)
	functions := kems[kemName]
	return Kem{id, functions}
}

func GetAllKems() []string {
	keys := make([]string, 0, len(kems))
	for k := range kems {
		keys = append(keys, k)
	}
	return keys
}

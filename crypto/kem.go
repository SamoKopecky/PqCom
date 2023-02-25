package crypto

import "github.com/SamoKopecky/pqcom/main/crypto/kem"

var kems = map[string]KemAlgorithm{
	"MyKyber":    &kem.MyKyber{},
	"CirclKyber": &kem.CirclKyber{},
}

type KemAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Dec(c, sk []byte) (key []byte)
	Enc(pk []byte) (c, key []byte)
	EkLen() (ekLen int)
}

type Kem struct {
	Id       int
	Functions KemAlgorithm
	EkLen     int
}

func GetKem(kemName string) Kem {
	id := getRowIndex(kems, kemName)
	functions := kems[kemName]
	return Kem{id, functions, functions.EkLen()}
}

func GetAllKems() []string {
	keys := make([]string, 0, len(kems))
	for k := range kems {
		keys = append(keys, k)
	}
	return keys
}

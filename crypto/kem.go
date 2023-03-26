package crypto

import (
	"github.com/SamoKopecky/pqcom/main/crypto/kem"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

func init() {
	var ids []uint8
	for name, alg := range kems {
		if slices.Contains(ids, alg.Id()) {
			log.Fatal().
				Int("id", int(alg.Id())).
				Str("name", name).
				Msg("Kem algorithm Id conflict, change to a different id")
		}
		ids = append(ids, alg.Id())
	}
}

var kems = map[string]KemAlgorithm{
	"PqComKyber512":  &kem.PqComKyber512{},
	"PqComKyber768":  &kem.PqComKyber768{},
	"PqComKyber1024": &kem.PqComKyber1024{},
	"CirclKyber512":  &kem.CirclKyber512{},
	"CirclKyber768":  &kem.CirclKyber768{},
	"CirclKyber1024": &kem.CirclKyber1024{},
}

type KemAlgorithm interface {
	KeyGen() (puK, prK []byte)
	Dec(c, prK []byte) (key []byte)
	Enc(puK []byte) (c, key []byte)
	EkLen() (ekLen int)
	CLen() (cLen int)
	Id() (id uint8)
}

type Kem struct {
	Id uint8
	F  KemAlgorithm
}

func GetKem(kemName string) Kem {
	functions := kems[kemName]
	return Kem{functions.Id(), functions}
}

func GetKemNames() []string {
	keys := make([]string, 0, len(kems))
	for k := range kems {
		keys = append(keys, k)
	}
	return keys
}

func GetKemById(id uint8) string {
	for k, v := range kems {
		if v.Id() == id {
			return k
		}
	}
	log.Fatal().Int("id", int(id)).Msg("Uknown kem id")
	return ""
}

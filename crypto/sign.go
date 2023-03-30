package crypto

import (
	"github.com/SamoKopecky/pqcom/main/crypto/sign"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

func init() {
	var ids []uint8
	for name, alg := range Signatures {
		if slices.Contains(ids, alg.Id()) {
			log.Fatal().
				Int("id", int(alg.Id())).
				Str("name", name).
				Msg("Sign algorithm Id conflict, change to a different id")
		}
		ids = append(ids, alg.Id())
	}
}

var Signatures = map[string]SignAlgorithm{
	"PqComDilithium2": &sign.PqComDilithium2{},
	"PqComDilithium3": &sign.PqComDilithium3{},
	"PqComDilithium5": &sign.PqComDilithium5{},
	"CirclDilithium2": &sign.CirclDilithium2{},
	"CirclDilithium3": &sign.CirclDilithium3{},
	"CirclDilithium5": &sign.CirclDilithium5{},
}

type SignAlgorithm interface {
	KeyGen() (puK, prK []byte)
	Verify(puK, msg, signature []byte) bool
	Sign(prK, msg []byte) (signature []byte)
	SignLen() (signLen int)
	PuKLen() (pkLen int)
	PrKLen() (skLen int)
	Id() (id uint8)
}

type Sign struct {
	Id uint8
	F  SignAlgorithm
}

func GetSign(signName string) Sign {
	functions := Signatures[signName]
	return Sign{functions.Id(), functions}
}

func GetSignNames() []string {
	keys := make([]string, 0, len(Signatures))
	for k := range Signatures {
		keys = append(keys, k)
	}
	return keys
}

func GetSignById(id uint8) string {
	for k, v := range Signatures {
		if v.Id() == id {
			return k
		}
	}
	log.Fatal().Int("id", int(id)).Msg("Uknown signature id")
	return ""
}

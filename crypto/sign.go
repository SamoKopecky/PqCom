package crypto

import (
	"github.com/SamoKopecky/pqcom/main/crypto/sign"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

func init() {
	var ids []uint8
	for name, alg := range signatures {
		if slices.Contains(ids, alg.Id()) {
			log.Fatal().
				Int("id", int(alg.Id())).
				Str("name", name).
				Msg("Sign algorithm Id conflict, change to a different id")
		}
		ids = append(ids, alg.Id())
	}
}

var signatures = map[string]SignAlgorithm{
	"PqComDilithium2": &sign.PqComDilithium2{},
	"CirclDilithium2": &sign.CirclDilithium2{},
}

type SignAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Verify(pk, msg, signature []byte) bool
	Sign(sk, msg []byte) (signature []byte)
	SignLen() (signLen int)
	PkLen() (pkLen int)
	SkLen() (skLen int)
	Id() (id uint8)
}

type Sign struct {
	Id uint8
	F  SignAlgorithm
}

func GetSign(signName string) Sign {
	functions := signatures[signName]
	return Sign{functions.Id(), functions}
}

func GetSignNames() []string {
	keys := make([]string, 0, len(signatures))
	for k := range signatures {
		keys = append(keys, k)
	}
	return keys
}

func GetSignById(id uint8) string {
	for k, v := range signatures {
		if v.Id() == id {
			return k
		}
	}
	log.Fatal().Int("id", int(id)).Msg("Uknown signature id")
	return ""
}

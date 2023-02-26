package crypto

import (
	"github.com/SamoKopecky/pqcom/main/crypto/sign"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

func init() {
	var ids []uint8
	for name, alg := range signatures {
		if slices.Contains(ids, alg.Id()) {
			log.WithFields(log.Fields{
				"id":   alg.Id(),
				"name": name,
			}).Fatal("Sign algorithm Id conflict, change to a different id")
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

func GetAllSigns() []string {
	keys := make([]string, 0, len(signatures))
	for k := range signatures {
		keys = append(keys, k)
	}
	return keys
}

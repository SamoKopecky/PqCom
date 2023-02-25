package crypto

import "github.com/SamoKopecky/pqcom/main/crypto/sign"

var signatures = map[string]SignAlgorithm{
	"MyDilithium": &sign.MyDilithium{},
	"CirclDilithium": &sign.CirclDilithium{},
}

type SignAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Verify(pk, msg, signature []byte) bool
	Sign(sk, msg []byte) (signature []byte)
	SignLen() (signLen int)
}

type Sign struct {
	Id        int
	Functions SignAlgorithm
}

func GetSign(signName string) Sign {
	id := getRowIndex(kems, signName)
	functions := signatures[signName]
	return Sign{id, functions}
}

func GetAllSigns() []string {
	keys := make([]string, 0, len(signatures))
	for k := range signatures {
		keys = append(keys, k)
	}
	return keys
}

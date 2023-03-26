package network

import (
	"fmt"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/rs/zerolog/log"
)

var (
	kem     crypto.Kem
	sign    crypto.Sign
	ekLen   int
	cLen    int
	signLen int
	sk      []byte
	pk      []byte
)

func SetupVars() {
	config := config.ReadConfig()

	kemStruct := crypto.GetKem(config.Kem)
	kem = kemStruct
	ekLen = kemStruct.F.EkLen()
	cLen = kem.F.CLen()

	signStruct := crypto.GetSign(config.Sign)
	sign = signStruct
	signLen = signStruct.F.SignLen()

	sk = config.Sk
	pk = config.Pk
}

type Stream struct {
	Conn      *net.TCPConn
	Msg       chan Msg
	key       []byte
	encrypt   bool
	aesCipher crypto.AesCipher
}

type Msg struct {
	Header Header
	Data   []byte
}

const CHUNK_SIZE = 2 << 13
const PACKET_SIZE = CHUNK_SIZE + HeaderLen + crypto.NONCE_LEN

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Error resolving")
		os.Exit(1)
	}
	return raddr
}

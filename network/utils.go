package network

import (
	"fmt"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/crypto"
	log "github.com/sirupsen/logrus"
)

var kem crypto.Kem
var sign crypto.Sign
var ekLen int
var signLen int
var sk []byte
var pk []byte

func SetupVars() {
	config := config.ReadConfig()

	kemStruct := crypto.GetKem(config.Kem)
	kem = kemStruct
	ekLen = kemStruct.F.EkLen()

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
const PACKET_SIZE = CHUNK_SIZE + HEADER_LEN

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error resolving")
		os.Exit(1)
	}
	return raddr
}

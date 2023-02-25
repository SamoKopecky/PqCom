package network

import (
	"fmt"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/crypto"
	log "github.com/sirupsen/logrus"
)

var kem crypto.KemAlgorithm
var sign crypto.SignAlgorithm
var ekLen int
var signLen int
var sk []byte
var pk []byte

func SetupVars() {
	config := config.ReadConfig()

	kemFuncs := crypto.GetKem(config.Kem).Functions
	kem = kemFuncs
	ekLen = kemFuncs.EkLen()

	signFncs := crypto.GetSign(config.Sign).Functions
	sign = signFncs
	signLen = signFncs.SignLen()

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

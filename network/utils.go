package network

import (
	"fmt"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/crypto"
	log "github.com/sirupsen/logrus"
)

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

const CHUNK_SIZE = 4096
const PACKET_SIZE = CHUNK_SIZE + HEADER_LEN

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error resolving")
		os.Exit(1)
	}
	return raddr
}

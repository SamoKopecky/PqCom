package network

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type Stream struct {
	Conn    *net.TCPConn
	Data    chan []byte
	key     []byte
	encrypt bool
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

package network

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type Stream struct {
	Conn *net.TCPConn
	Data chan []byte
}

const ChunkSize = 4 << 12

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error resolving")
		os.Exit(1)
	}
	return raddr
}

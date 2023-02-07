package network

import (
	"fmt"
	"net"
	"os"
)

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error resolving")
		os.Exit(1)
	}
	return raddr
}

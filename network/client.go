package network

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type client struct {
	sock *net.TCPConn
}

func (c *client) connect(addr string, port int) {
	prot := "tcp"
	var err error
	c.sock, err = net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error trying to connect")
		os.Exit(1)
	}
	log.WithField("addr", c.sock.RemoteAddr()).Info("Connected")
}

func (c *client) send(data []byte) {
	n, err := c.sock.Write(data)
	log.WithField("len", n).Debug("Send data to socket")
	if err != nil {
		log.WithField("error", err).Error("Can't write to socket")
		os.Exit(1)
	}
}

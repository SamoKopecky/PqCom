package network

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	Sock *net.TCPConn
}

func (c *Client) Connect(addr string, port int) {
	prot := "tcp"
	var err error
	c.Sock, err = net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error trying to connect")
		os.Exit(1)
	}
	log.WithField("addr", c.Sock.RemoteAddr()).Info("Connected")
}

func (c *Client) Send(data []byte) {
	n, err := c.Sock.Write(data)
	log.WithField("len", n).Debug("Send data to socket")
	if err != nil {
		log.WithField("error", err).Error("Can't write to socket")
		os.Exit(1)
	}
}

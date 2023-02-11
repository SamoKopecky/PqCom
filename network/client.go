package network

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type client struct {
	sock *net.TCPConn
	send chan []byte
}

func (c *client) connect(addr string, port int) {
	prot := "tcp"
	var err error
	c.sock, err = net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error trying to connect")
		os.Exit(1)
	}
	log.WithField("conn", c.sock.RemoteAddr()).Info("Connected")
}

func (c *client) startSend() {
	for {
		data := <-c.send
		_, err := c.sock.Write(data)
		if err != nil {
			log.WithField("error", err).Error("Can't write to socket")
			os.Exit(1)
		}
	}
}

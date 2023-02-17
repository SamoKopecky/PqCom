package network

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func Connect(addr string, port int) Stream {
	prot := "tcp"
	conn, err := net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	stream := Stream{Data: make(chan []byte)}
	if err != nil {
		log.WithField("error", err).Error("Error trying to connect")
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"remote addr": conn.RemoteAddr(),
		"local addr":  conn.LocalAddr(),
	}).Info("Connected")
	stream.Conn = conn
	go stream.readData()
	return stream
}

func (s *Stream) Send(data []byte) {
	n, err := s.Conn.Write(data)
	log.WithField("len", n).Debug("Send data to socket")
	if err != nil {
		log.WithField("error", err).Error("Can't write to socket")
		os.Exit(1)
	}
}

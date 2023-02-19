package network

import (
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/kyber"
	log "github.com/sirupsen/logrus"
)

func Connect(addr string, port int) Stream {
	prot := "tcp"
	conn, err := net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	s := Stream{Data: make(chan []byte)}
	if err != nil {
		log.WithField("error", err).Error("Error trying to connect")
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"remote addr": conn.RemoteAddr(),
		"local addr":  conn.LocalAddr(),
	}).Info("Connected")
	s.Conn = conn

	s.clientKeyEnc()
	go s.readData()
	return s
}

func (s *Stream) clientKeyEnc() {
	pk, sk := kyber.CcakemKeyGen()
	s.Send(s.packWithHeader(pk))
	c := s.readWithHeader()
	key := kyber.CcakemDec(c, sk)
	s.key = key
	// fmt.Printf("%d\n", s.key)
}

func (s *Stream) packWithHeader(data []byte) (dataWithHeader []byte) {
	header := Header{Len: uint16(len(data))}
	dataWithHeader = append(header.build(), data...)
	return
}

func (s *Stream) Send(data []byte) {
	n, err := s.Conn.Write(data)
	log.WithField("len", n).Debug("Send data to socket")
	if err != nil {
		log.WithField("error", err).Fatal("Can't write to socket")
	}
}

func (s *Stream) SafeSend(data []byte) {
	s.Send(s.packWithHeader(data))
}

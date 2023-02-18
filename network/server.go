package network

import (
	"fmt"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/io"
	"github.com/SamoKopecky/pqcom/main/kyber"
	log "github.com/sirupsen/logrus"
)

func Listen(port int, streamFactory chan<- Stream, always bool) {
	prot := "tcp"
	address := resolvedAddr(prot, "0.0.0.0", port)
	listener, err := net.ListenTCP(prot, address)

	if err != nil {
		log.WithField("error", err).Error("Error trying to listen")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		log.WithField("addr", address).Info("Acepting connections")
		s := Stream{Data: make(chan []byte)}
		s.Conn, err = listener.AcceptTCP()
		if err != nil {
			log.WithField("error", err).Error("Error accpeting conn")
			s.Conn.Close()
		}
		log.WithField("remote addr", s.Conn.RemoteAddr()).Info("Recevied connection")

		s.initServer()
		streamFactory <- s
		go s.readData()
		if !always {
			break
		}
	}
}

func (s *Stream) initServer() {
	pk := s.readWithHeader()
	c, key := kyber.CcakemEnc(pk)
	s.Send(s.packWithHeader(c))
	fmt.Printf("%d\n", key)
}

func (s *Stream) readWithHeader() (data []byte) {
	var readBytes int
	var recentData []byte
	var header Header
	var n int

	buf := make([]byte, 0, ChunkSize)
	n = io.Read(s.Conn, buf)
	buf = buf[:n]
	header.parse(buf)
	recentData = buf[HEADER_LEN:n]

	for {
		data = append(data, recentData...)
		readBytes += len(recentData)
		if readBytes >= int(header.Len) {
			break
		}
		n = io.Read(s.Conn, buf)
		recentData = buf[:n]
	}
	return
}

func (s *Stream) readData() {
	defer s.Conn.Close()
	io.ReadByChunks(s.Conn, s.Data, ChunkSize)
	log.WithField("remote addr", s.Conn.RemoteAddr()).Info("Connection ended, closing")
}

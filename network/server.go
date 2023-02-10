package network

import (
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type server struct {
	conn     *net.TCPConn
	listener *net.TCPListener
	c        chan string
}

func (s *server) listen(port int) {
	prot := "tcp"
	address := resolvedAddr(prot, "0.0.0.0", port)
	var err error
	s.listener, err = net.ListenTCP(prot, address)

	if err != nil {
		log.WithField("error", err).Error("Error trying to listen")
		os.Exit(1)
	}
	defer s.listener.Close()
	log.WithField("addr", address).Info("Listening")

	for {
		s.conn, err = s.listener.AcceptTCP()
		if err != nil {
			log.WithField("error", err).Error("Error accpeting conn")
			s.conn.Close()
		}
		log.WithField("remote addr", s.conn.RemoteAddr()).Info("Connected")

		go s.handleConnection()
	}

}

func (s *server) handleConnection() {
	defer s.conn.Close() // clean up when done
	var buf []byte

	for {
		buf = make([]byte, 2048)
		_, err := s.conn.Read(buf)

		if err != nil {
			log.WithField("error", err).Error("Error reading from accpeted conn")
			return
		}
		s.c <- string(buf)
	}
}

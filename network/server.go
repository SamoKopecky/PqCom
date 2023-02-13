package network

import (
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/handler"
	log "github.com/sirupsen/logrus"
)

type server struct {
	conn     *net.TCPConn
	listener *net.TCPListener
	recv     chan []byte
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
		log.WithField("remote addr", s.conn.RemoteAddr()).Info("Recevied connection")

		go s.handleConnection()
	}

}

func (s *server) handleConnection() {
	defer s.conn.Close()
	handler.ReadByChunks(s.conn, s.recv, chunkSize)
	log.WithField("remote addr", s.conn.RemoteAddr()).Info("Connection ended, closing")
}

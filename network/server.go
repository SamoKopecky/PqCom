package network

import (
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/io"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Listener *net.TCPListener
}

type Stream struct {
	Conn *net.TCPConn
	Data chan []byte
}

func (s *Server) Listen(port int, streamFactory chan<- Stream) {
	prot := "tcp"
	address := resolvedAddr(prot, "0.0.0.0", port)
	var err error
	s.Listener, err = net.ListenTCP(prot, address)

	if err != nil {
		log.WithField("error", err).Error("Error trying to listen")
		os.Exit(1)
	}
	defer s.Listener.Close()
	log.WithField("addr", address).Info("Listening")

	for {
		stream := Stream{Data: make(chan []byte)}
		stream.Conn, err = s.Listener.AcceptTCP()
		if err != nil {
			log.WithField("error", err).Error("Error accpeting conn")
			stream.Conn.Close()
		}
		log.WithField("remote addr", stream.Conn.RemoteAddr()).Info("Recevied connection")

		streamFactory <- stream
		go stream.handleConnection()
	}

}

func (stream *Stream) handleConnection() {
	defer stream.Conn.Close()
	io.ReadByChunks(stream.Conn, stream.Data, ChunkSize)
	log.WithField("remote addr", stream.Conn.RemoteAddr()).Info("Connection ended, closing")
}

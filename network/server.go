package network

import (
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/io"
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
		stream := Stream{Data: make(chan []byte)}
		stream.Conn, err = listener.AcceptTCP()
		if err != nil {
			log.WithField("error", err).Error("Error accpeting conn")
			stream.Conn.Close()
		}
		log.WithField("remote addr", stream.Conn.RemoteAddr()).Info("Recevied connection")

		streamFactory <- stream
		go stream.readData()
		if !always {
			break
		}
	}
}

func (stream *Stream) readData() {
	defer stream.Conn.Close()
	io.ReadByChunks(stream.Conn, stream.Data, ChunkSize)
	log.WithField("remote addr", stream.Conn.RemoteAddr()).Info("Connection ended, closing")
}

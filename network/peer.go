package network

import (
	"io"
	"os"

	"github.com/SamoKopecky/pqcom/main/handler"
	log "github.com/sirupsen/logrus"
)

var p *peer

type peer struct {
	c *client
	s *server
}

func init() {
	p = &peer{&client{}, &server{recv: make(chan []byte)}}
}

func Chat(destAddr string, srcPort, destPort int) {
	go p.s.listen(srcPort)
	go handler.Printer(p.s.recv, false)

	handler.ReadUserInput("Press enter to connect\n")
	p.c.connect(destAddr, destPort)
	defer p.c.sock.Close()
	for {
		data := []byte(handler.ReadUserInput(""))
		p.c.send(data)
	}
}

func Send(destAddr string, srcPort, destPort int, filePath string) {
	p.c.connect(destAddr, destPort)
	defer p.c.sock.Close()
	chunks := make(chan []byte)
	var source io.Reader
	var err error

	if filePath != "" {
		source, err = os.Open(filePath)
		if err != nil {
			log.WithField("error", err).Error("Error opening file")
		}
	} else {
		source = os.Stdin
	}
	go func() {
		handler.ReadByChunks(source, chunks, chunkSize)
		close(chunks)
	}()
	for msg := range chunks {
		p.c.send(msg)
	}
}

func Receive(destAddr string, srcPort, destPort int, dir string) {
	if dir != "" {
		go handler.FileWriter(dir, p.s.recv)
	} else {
		go handler.Printer(p.s.recv, true)
	}
	p.s.listen(srcPort)
}

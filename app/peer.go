package app

import (
	"io"
	"os"

	myio "github.com/SamoKopecky/pqcom/main/io"
	"github.com/SamoKopecky/pqcom/main/network"
	log "github.com/sirupsen/logrus"
)

var p *peer

type peer struct {
	c *network.Client
	s *network.Server
}

func init() {
	p = &peer{&network.Client{}, &network.Server{}}
}

func Chat(destAddr string, srcPort, destPort int) {
	streamFactory := make(chan network.Stream)
	go p.s.Listen(srcPort, streamFactory)
	go func() {
		conn := <-streamFactory
		go printer(conn, false)
	}()

	myio.ReadUserInput("Press enter to connect\n")
	p.c.Connect(destAddr, destPort)
	defer p.c.Sock.Close()
	for {
		data := []byte(myio.ReadUserInput(""))
		p.c.Send(data)
	}
}

func Send(destAddr string, srcPort, destPort int, filePath string) {
	p.c.Connect(destAddr, destPort)
	defer p.c.Sock.Close()
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
		myio.ReadByChunks(source, chunks, network.ChunkSize)
		close(chunks)
	}()
	for msg := range chunks {
		p.c.Send(msg)
	}
}

func Receive(destAddr string, srcPort, destPort int, dir string) {
	streamFactory := make(chan network.Stream)
	go p.s.Listen(srcPort, streamFactory)
	for {
		stream := <-streamFactory
		if dir != "" {
			go dirFileWriter(stream.Data, dir)
		} else {
			go printer(stream, true)
		}
	}
}

package network

import (
	"fmt"
	"io"
	"os"

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

func (s *server) printer(clean bool) {
	for {
		msg := <-s.recv
		if clean {
			fmt.Printf("%s", string(msg))
			continue
		}
		fmt.Printf("[%s]: %s", s.conn.RemoteAddr(), string(msg))
	}
}

func Chat(destAddr string, srcPort, destPort int) {
	go p.s.listen(srcPort)
	go p.s.printer(false)

	readUserInput("Press enter to connect\n")
	p.c.connect(destAddr, destPort)
	defer p.c.sock.Close()
	for {
		data := []byte(readUserInput(""))
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
	go readByChunks(source, chunks)
	for msg := range chunks {
		p.c.send(msg)
	}
}

func Receive(destAddr string, srcPort, destPort int, dir string) {
	if dir != "" {
		print("TODO\n")
		os.Exit(0)
	} else {
		go p.s.printer(true)
	}
	p.s.listen(srcPort)
}

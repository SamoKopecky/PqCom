package network

import (
	"fmt"
	"os"
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
	if filePath != "" {
		print("TODO\n")
		os.Exit(0)
	} else {
		p.c.connect(destAddr, destPort)
		data := readStdin()
		p.c.send(data)
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

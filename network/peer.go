package network

import (
	"bufio"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

const randomBytes = 5

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
		go p.s.fileWriter(dir)
	} else {
		go p.s.printer(true)
	}
	p.s.listen(srcPort)
}

func (s *server) fileWriter(dir string) {
	newFile := true
	var file *os.File
	var err error

	for {
		msg := <-s.recv

		if newFile {
			// TODO: Do this nicer
			fileName := fmt.Sprintf("pqcom_temp_%s", RandStringBytes(randomBytes))
			for containsDir(fileName, dir) {
				fileName = fmt.Sprintf("pqcom_temp_%s", RandStringBytes(randomBytes))
			}
			filePath := fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileName)

			file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.WithField("error", err).Error("Error opening file")
			}
			defer file.Close()
			newFile = false
		}

		w := bufio.NewWriter(file)
		n, err := w.Write(msg)
		if err != nil {
			panic(err)
		}
		if n == 0 {
			newFile = true
		}
		w.Flush()
	}
}

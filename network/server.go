package network

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func listen(addr string, port int) {
	prot := "tcp"
	address := resolvedAddr(prot, addr, port)
	listen, err := net.ListenTCP(prot, address)
	if err != nil {
		log.WithField("error", err).Error("Error trying to listen")
		os.Exit(1)
	}
	defer listen.Close()
	log.WithField("addr", address).Info("Listening")

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.WithField("error", err).Error("Error accpeting conn")
			conn.Close()
			continue
		}
		log.WithField("remote addr", conn.RemoteAddr()).Info("Connected")

		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close() // clean up when done
	var buf []byte

	for {
		buf = make([]byte, 2048)
		n, err := conn.Read(buf)

		if err != nil {
			log.WithField("error", err).Error("Error reading from accpeted conn")
			return
		}

		fmt.Print("I read: ", string(buf))

		write, err := conn.Write(buf[:n])
		if err != nil {
			fmt.Println("failed to write to client:", err)
			return
		}
		if write != n { // was all data sent
			fmt.Println("warning: not all data sent to client")
			return
		}
	}
}

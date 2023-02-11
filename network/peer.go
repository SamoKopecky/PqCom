package network

import (
	"fmt"
)

func Start(destAddr string, srcPort, destPort int, stdin bool) {
	client := client{send: make(chan []byte)}
	server := server{recv: make(chan []byte)}

	go server.listen(srcPort)
	go server.printer()

	readUserInput("Press enter to connect\n")
	client.connect(destAddr, destPort)
	go client.startSend()

	if stdin {
		// data := readStdin()
		// fmt.Printf("%s", data)
		// send(daddr, dport, []byte(data))
	} else {
		for {
			data := []byte(readUserInput(""))
			client.send <- data
		}
	}
}

func (s *server) printer() {
	for {
		msg := <-s.recv
		fmt.Printf("[%s]: %s", s.conn.RemoteAddr(), string(msg))
	}
}

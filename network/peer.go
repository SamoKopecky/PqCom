package network

import (
	"fmt"
)

func Start(destAddr string, srcPort, destPort int, stdin bool) {
	client := client{send: make(chan []byte)}
	server := server{recv: make(chan string)}

	go server.listen(srcPort)
	go printer(server.recv)

	readInput("Press anything to connect\n")
	client.connect(destAddr, destPort)
	go client.startSend()

	if stdin {
		// data := readStdin()
		// fmt.Printf("%s", data)
		// send(daddr, dport, []byte(data))
	} else {
		for {
			data := []byte(readInput("\nme: "))
			client.send <- data
		}
	}
}

func printer(c <-chan string) {
	for {
		msg := <-c
		fmt.Printf("\nother: %s\n", msg)
	}
}

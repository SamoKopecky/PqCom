package network

import (
	"fmt"
)

func Start(daddr string, lport, dport int, stdin bool) {
	client := new(client)
	server := new(server)
	server.c = make(chan string)

	go server.listen(lport)
	go printer(server.c)

	readInput("Press anything to connect\n")
	client.connect(daddr, dport)

	if stdin {
		// data := readStdin()
		// fmt.Printf("%s", data)
		// send(daddr, dport, []byte(data))
	} else {
		for {
			data := []byte(readInput("\nme: "))
			client.send(data)
		}
	}
}

func printer(c <-chan string) {
	for {
		msg := <-c
		fmt.Printf("\nother: %s\n", msg)
	}
}

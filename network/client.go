package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func send(addr string, port int) {
	prot := "tcp"
	conn, err := net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error trying to connect")
		os.Exit(1)
	}
	defer conn.Close()

	for {
		fmt.Print("Type something to send: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// read response
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("failed reading response:", err)
			os.Exit(1)
		}
		fmt.Print("Response: ", string(buf[:n]))
	}
}

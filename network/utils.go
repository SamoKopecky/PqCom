package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error resolving")
		os.Exit(1)
	}
	return raddr
}

func readStdin() []byte {
	var data []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return data
}

func readInput(message string) string {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

package network

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
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
	r := bufio.NewReader(os.Stdin)
	var content []byte
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		content = append(content, buf...)
		// process buf
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
	return content
}

func readUserInput(promt string) string {
	fmt.Print(promt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

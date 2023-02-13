package network

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const chunkSize = 4 << 10

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func containsDir(file string, dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		if file == e.Name() {
			return true
		}
	}
	return false
}

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error resolving")
		os.Exit(1)
	}
	return raddr
}

func readByChunks(reader io.Reader, chunks chan<- []byte) {
	r := bufio.NewReader(reader)
	for {
		// TODO: Make is that buf doesn't have to initialize every time
		buf := make([]byte, 0, chunkSize)
		n, err := r.Read(buf[:cap(buf)])
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		chunks <- buf[:n]
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}
	close(chunks)
}

func readUserInput(promt string) string {
	fmt.Print(promt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

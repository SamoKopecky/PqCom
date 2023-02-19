package network

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type Stream struct {
	Conn *net.TCPConn
	Data chan []byte
	key  []byte
}

const CHUNK_SIZE = 2 << 11
const PACKET_SIZE = CHUNK_SIZE + HEADER_LEN

func resolvedAddr(prot string, addr string, port int) *net.TCPAddr {
	raddr, err := net.ResolveTCPAddr(prot, fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.WithField("error", err).Error("Error resolving")
		os.Exit(1)
	}
	return raddr
}

// var readBytes int
// var recentData []byte
// var header Header
// var n int
// var err error
// var dataBuf []byte

// buf := make([]byte, 0, CHUNK_SIZE+HEADER_LEN)
// readBytes = 0
// for {
// 	n, err = myio.Read(s.Conn, buf)
// 	if err == io.EOF && n == 0 {
// 		break
// 	}
// 	buf = buf[:n]
// 	header.parse(buf)
// 	recentData = buf[HEADER_LEN:n]

// 	for {
// 		readBytes += len(recentData)
// 		if readBytes >= int(header.Len) {
// 			if readBytes != int(header.Len) {
// 				readBytes = readBytes - int(header.Len)
// 				dataBuf = append(dataBuf, recentData[:readBytes]...)
// 			} else {
// 				readBytes = 0
// 				dataBuf = append(dataBuf, recentData...)
// 			}
// 			// Decrypt
// 			s.Data <- dataBuf
// 			break
// 		}
// 		dataBuf = append(dataBuf, recentData...)
// 		n, _ = myio.Read(s.Conn, buf)
// 		recentData = buf[:n]
// 	}
// }

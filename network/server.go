package network

import (
	"io"
	"net"
	"os"

	myio "github.com/SamoKopecky/pqcom/main/io"
	"github.com/SamoKopecky/pqcom/main/kyber"
	log "github.com/sirupsen/logrus"
)

func Listen(port int, streamFactory chan<- Stream, always bool) {
	prot := "tcp"
	address := resolvedAddr(prot, "0.0.0.0", port)
	listener, err := net.ListenTCP(prot, address)

	if err != nil {
		log.WithField("error", err).Error("Error trying to listen")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		log.WithField("addr", address).Info("Acepting connections")
		s := Stream{Data: make(chan []byte)}
		s.Conn, err = listener.AcceptTCP()
		if err != nil {
			log.WithField("error", err).Error("Error accpeting conn")
			s.Conn.Close()
		}
		log.WithField("remote addr", s.Conn.RemoteAddr()).Info("Recevied connection")

		s.serverKeyEnc()
		streamFactory <- s
		go s.readData()
		if !always {
			break
		}
	}
}

func (s *Stream) serverKeyEnc() {
	pk := s.readWithHeader()
	c, key := kyber.CcakemEnc(pk)
	s.Send(s.packWithHeader(c))
	s.key = key
	// fmt.Printf("%d\n", s.key)
}

func (s *Stream) readWithHeader() (data []byte) {
	var readBytes int
	var recentData []byte
	var header Header
	var n int

	buf := make([]byte, 0, CHUNK_SIZE+HEADER_LEN)
	n, _ = myio.Read(s.Conn, buf)
	buf = buf[:n]
	header.parse(buf)
	recentData = buf[HEADER_LEN:n]

	for {
		data = append(data, recentData...)
		readBytes += len(recentData)
		if readBytes >= int(header.Len) {
			break
		}
		n, _ = myio.Read(s.Conn, buf)
		recentData = buf[:n]
	}
	return
}

func (s *Stream) readWithHeaderChan() {
	buf := make([]byte, 0, PACKET_SIZE)
	var n int
	var err error
	var header Header
	var packetRead int
	var packetData []byte
	var leftToRead int
	var first bool
	lastPacketEnd := 0

	for {
		if lastPacketEnd == 0 {
			n, err = myio.Read(s.Conn, buf)
			if err == io.EOF && n == 0 {
				break
			}
		}
		if n-lastPacketEnd < HEADER_LEN {
			temp := []byte{}
			temp = append(temp, buf[lastPacketEnd:]...)
			n, err = myio.Read(s.Conn, buf)
			buf = buf[:n]
			temp = append(temp, buf[:n]...)
			header.parse(temp[:HEADER_LEN])
			lastPacketEnd = -(n - lastPacketEnd)
		} else {
			header.parse(buf[lastPacketEnd:n])
		}
		packetData = append([]byte{}, buf[lastPacketEnd+HEADER_LEN:n]...)
		packetRead += len(packetData)

		first = true

		for {
			if packetRead >= int(header.Len) {
				if packetRead == int(header.Len) {
					lastPacketEnd = 0
					if !first {
						packetData = append(packetData, buf[:n]...)
					} else {
						first = false
					}
				} else {
					leftToRead = int(header.Len) - (packetRead - n)
					packetData = append(packetData, buf[:leftToRead]...)
					lastPacketEnd = leftToRead
				}
				// fmt.Printf("%s", string(packetData))
				s.Data <- append([]byte{}, packetData...)
				packetRead = 0
				break
			} else {
				if !first {
					packetData = append(packetData, buf[:n]...)
				} else {
					first = false
				}
			}
			n, err = myio.Read(s.Conn, buf)
			buf = buf[:n]
			packetRead += n
		}
	}
}

func (s *Stream) readData() {
	defer s.Conn.Close()
	s.readWithHeaderChan()
	log.WithField("remote addr", s.Conn.RemoteAddr()).Info("Connection ended, closing")
}

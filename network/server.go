package network

import (
	"io"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/crypto"
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
		s := Stream{Data: make(chan []byte), encrypt: false}
		s.Conn, err = listener.AcceptTCP()
		if err != nil {
			log.WithField("error", err).Error("Error accpeting conn")
			s.Conn.Close()
		}
		log.WithField("remote addr", s.Conn.RemoteAddr()).Info("Recevied connection")

		s.aesCipher = crypto.AesCipher{}
		s.serverKeyEnc()
		streamFactory <- s
		go s.readData()
		if !always {
			break
		}
	}
}

func (s *Stream) serverKeyEnc() {
	pk := s.readPacket()
	c, key := kyber.CcakemEnc(pk)
	s.Send(c)
	nonce := s.readPacket()
	s.key = key
	s.encrypt = true
	s.aesCipher.Create(s.key, nonce)
}

func (s *Stream) readPacket() (data []byte) {
	go func() {
		s.readChunks(true)
		close(s.Data)
	}()
	for msg := range s.Data {
		data = append(data, msg...)
	}
	s.Data = make(chan []byte)
	return data
}

func (s *Stream) readChunks(onePacket bool) {
	buf := make([]byte, 0, PACKET_SIZE)
	var err error
	var header Header
	var packetRead, lastPacketEnd, n, pack int
	var packetData []byte
	var first bool

	for {
		if lastPacketEnd == 0 {
			n, err = myio.Read(s.Conn, buf)
			buf = buf[:n]
			if err == io.EOF && n == 0 {
				break
			}
		}
		// if n-lastPacketEnd < HEADER_LEN {
		// 	temp = append([]byte{}, buf[lastPacketEnd:]...)
		// 	n, err = myio.Read(s.Conn, buf)
		// 	buf = buf[:n]
		// 	temp = append(temp, buf...)
		// 	lastPacketEnd = -(n - lastPacketEnd)
		// 	header.parse(temp[:HEADER_LEN])
		// } else {
		// }
		header.parse(buf[lastPacketEnd:n])
		packetData = append([]byte{}, buf[lastPacketEnd+HEADER_LEN:n]...)
		packetRead += len(packetData)

		first = true
		for {
			if packetRead == int(header.Len) {
				lastPacketEnd = 0
				pack = n
			} else if packetRead > int(header.Len) {
				lastPacketEnd = int(header.Len) - (packetRead - n)
				pack = lastPacketEnd
			} else {
				pack = n
			}
			if !first {
				packetData = append(packetData, buf[:pack]...)
			}
			if packetRead >= int(header.Len) {
				break
			}
			n, err = myio.Read(s.Conn, buf)
			buf = buf[:n]
			packetRead += n
			first = false
		}
		if s.encrypt {
			packetData = s.aesCipher.Decrypt(packetData)
		}
		s.Data <- append([]byte{}, packetData...)
		packetRead = 0
		if onePacket {
			return
		}
	}
}

func (s *Stream) readData() {
	defer s.Conn.Close()
	s.readChunks(false)
	log.WithField("remote addr", s.Conn.RemoteAddr()).Info("Connection ended, closing")
}

package network

import (
	"io"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/crypto"
	myio "github.com/SamoKopecky/pqcom/main/io"
	log "github.com/sirupsen/logrus"
)

func Listen(port int, streamFactory chan<- Stream, always bool) {
	SetupVars()
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
		s := Stream{Msg: make(chan Msg), encrypt: false}
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
	clientInit := ClientInit{}
	signedData := clientInit.parse(s.readPacket())
	if clientInit.kemType != kem.Id || clientInit.signType != sign.Id {
		errorReason := "Config algorithm mismtatch"
		errorMsg := ErrorMsg{errorReason}
		s.Send(errorMsg.build(), ErrorT)
		log.WithFields(log.Fields{
			"kem id":           kem.Id,
			"received kem id":  clientInit.kemType,
			"sign id":          sign.Id,
			"received sign id": clientInit.signType,
		}).Fatal(errorReason)
	}

	nonce := clientInit.nonce
	signature := clientInit.sig
	if !sign.F.Verify(pk, signedData, signature) {
		log.Fatal("Signaute failure")
	}
	c, key := kem.F.Enc(myio.Copy(clientInit.eK))

	serverInit := ServerInit{keyC: c}
	s.Send(serverInit.build(), ServerInitT)

	s.key = key
	s.encrypt = true
	s.aesCipher.Create(s.key, nonce)
}

func (s *Stream) readPacket() (data []byte) {
	go func() {
		s.readChunks()
		close(s.Msg)
	}()
	for msg := range s.Msg {
		if msg.Header.Type == ErrorT {
			errorMsg := ErrorMsg{}
			errorMsg.parse(msg.Data)
			log.WithField("error", errorMsg.reason).Fatal("Received error from other peer")
		}
		data = append(data, msg.Data...)
	}
	s.Msg = make(chan Msg)
	return data
}

func (s *Stream) readChunks() {
	buf := make([]byte, 0, PACKET_SIZE)
	var err error
	var packetRead, lastPacketEnd, n, pack int
	var packetData []byte
	var first bool
	var msg Msg

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
		msg.Header.parse(buf[lastPacketEnd:n])
		packetData = append([]byte{}, buf[lastPacketEnd+HEADER_LEN:n]...)
		packetRead += len(packetData)

		first = true
		for {
			if packetRead == int(msg.Header.Len) {
				lastPacketEnd = 0
				pack = n
			} else if packetRead > int(msg.Header.Len) {
				lastPacketEnd = int(msg.Header.Len) - (packetRead - n)
				pack = lastPacketEnd
			} else {
				pack = n
			}
			if !first {
				packetData = append(packetData, buf[:pack]...)
			}
			if packetRead >= int(msg.Header.Len) {
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
		msg.Data = myio.Copy(packetData)
		s.Msg <- msg
		packetRead = 0
		if msg.Header.Type != ContentT {
			return
		}
	}
}

func (s *Stream) readData() {
	defer s.Conn.Close()
	s.readChunks()
	log.WithField("remote addr", s.Conn.RemoteAddr()).Info("Connection ended, closing")
}

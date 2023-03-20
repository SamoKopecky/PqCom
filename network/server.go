package network

import (
	"io"
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/cookie"
	"github.com/SamoKopecky/pqcom/main/crypto"
	myio "github.com/SamoKopecky/pqcom/main/myio"
	"github.com/rs/zerolog/log"
)

func Listen(port int, streamFactory chan<- Stream, always bool) {
	SetupVars()
	prot := "tcp"
	address := resolvedAddr(prot, "0.0.0.0", port)
	listener, err := net.ListenTCP(prot, address)

	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Error trying to listen")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		log.Info().Str("addr", address.String()).Msg("Acepting connections")
		s := Stream{Msg: make(chan Msg), encrypt: false}
		s.Conn, err = listener.AcceptTCP()
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("Error accpeting conn")
			s.Conn.Close()
		}
		log.Info().Str("remote addr", s.Conn.RemoteAddr().String()).Msg("Recevied connection")

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
	log.Info().Msg("starting server key encapsulation")

	ci := ClientInit{}
	ci.parse(s.readPacket())
	payload := ci.payload()
	nonce := ci.nonce
	signature := ci.sig
	log.Debug().Msg("Verifing signature")
	if !sign.F.Verify(pk, payload, signature) {
		log.Fatal().Msg("Signature verification failed")
	}

	cookie := cookie.Cookie{Seed: pk, Timestamp: ci.timestamp}
	if !cookie.Exists() {
		cookie.Save()
	}
	if !cookie.IsNewer() {
		errorReason := "Old timestamp"
		errorMsg := ErrorMsg{errorReason}
		s.Send(errorMsg.build(), ErrorT)
		log.Fatal().Msg(errorReason)
	}

	if ci.kemType != kem.Id || ci.signType != sign.Id {
		errorReason := "Config algorithm mismtatch"
		errorMsg := ErrorMsg{errorReason}
		s.Send(errorMsg.build(), ErrorT)
		log.Fatal().
			Int("kem id", int(kem.Id)).
			Int("received kem id", int(ci.kemType)).
			Int("sign id", int(sign.Id)).
			Int("received sign id", int(ci.signType)).
			Msg(errorReason)
	}

	log.Debug().Msg("Encapsulating shared key")
	c, key := kem.F.Enc(myio.Copy(ci.eK))

	signature = sign.F.Sign(sk, c)
	si := ServerInit{keyC: c, sig: signature}
	s.Send(si.build(), ServerInitT)
	cookie.Save()
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
			log.Fatal().Str("error", errorMsg.reason).Msg("Received error from other peer")
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
			log.Debug().Int("len", len(packetData)).Msg("Decrypting data")
			packetData = s.aesCipher.Decrypt(packetData)
		}
		msg.Data = myio.Copy(packetData)
		s.Msg <- msg
		packetRead = 0
		log.Info().Int("msg type", int(msg.Header.Type)).Msg("Received msg")
		if msg.Header.Type != ContentT {
			return
		}
	}
}

func (s *Stream) readData() {
	defer s.Conn.Close()
	s.readChunks()
	log.Info().Str("remote addr", s.Conn.RemoteAddr().String()).Msg("Connection ended, closing")
}

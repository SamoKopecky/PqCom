package network

import (
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/crypto"
	log "github.com/sirupsen/logrus"
)

func Connect(addr string, port int) Stream {
	SetupVars()
	prot := "tcp"
	conn, err := net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	s := Stream{Msg: make(chan Msg), encrypt: false}
	if err != nil {
		log.WithField("error", err).Error("Error trying to connect")
		os.Exit(1)
	}
	log.WithFields(log.Fields{
		"remote addr": conn.RemoteAddr(),
		"local addr":  conn.LocalAddr(),
	}).Info("Connected")
	s.Conn = conn
	s.aesCipher = crypto.AesCipher{}
	s.clientKeyEnc()
	go s.readData()
	return s
}

func (s *Stream) clientKeyEnc() {
	ek, dk := kem.KeyGen()
	nonce := crypto.GenerateNonce()
	clientInit := ClientInit{eK: ek, nonce: nonce}
	payload := clientInit.build()
	signature := sign.Sign(sk, payload)

	clientInit.sig = signature
	s.Send(clientInit.build(), ClientInitT)

	serverInit := ServerInit{}
	serverInit.parse(s.readPacket())

	key := kem.Dec(serverInit.keyC, dk)
	s.key = key
	s.aesCipher.Create(s.key, nonce)
	s.encrypt = true
}

func (s *Stream) Send(data []byte, dataType Type) {
	if s.encrypt {
		data = s.aesCipher.Encrypt(data)
	}
	header := Header{Len: uint16(len(data)), Type: dataType}
	n, err := s.Conn.Write(append(header.build(), data...))
	log.WithField("len", n).Debug("Send data to socket")
	if err != nil {
		log.WithField("error", err).Fatal("Can't write to socket")
	}
}

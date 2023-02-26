package network

import (
	"net"
	"os"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/rs/zerolog/log"
)

func Connect(addr string, port int) Stream {
	SetupVars()
	prot := "tcp"
	conn, err := net.DialTCP(prot, nil, resolvedAddr(prot, addr, port))
	s := Stream{Msg: make(chan Msg), encrypt: false}
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Error trying to connect")
		os.Exit(1)
	}
	log.Info().
		Str("remote addr", conn.RemoteAddr().String()).
		Str("local addr", conn.LocalAddr().String()).
		Msg("Connected")
	s.Conn = conn
	s.aesCipher = crypto.AesCipher{}
	s.clientKeyEnc()
	go s.readData()
	return s
}

func (s *Stream) clientKeyEnc() {
	log.Info().Msg("starting client key encapsulation")
	ek, dk := kem.F.KeyGen()
	nonce := crypto.GenerateNonce()
	clientInit := ClientInit{
		eK:       ek,
		nonce:    nonce,
		kemType:  kem.Id,
		signType: sign.Id,
	}
	payload := clientInit.build()
	log.Debug().Msg("Signing payload")

	clientInit.sig = sign.F.Sign(sk, payload)
	log.Debug().Msg("Sending client init")
	s.Send(clientInit.build(), ClientInitT)

	serverInit := ServerInit{}
	serverInit.parse(s.readPacket())

	log.Debug().Msg("Decapsulating shared key")
	key := kem.F.Dec(serverInit.keyC, dk)
	s.key = key
	s.aesCipher.Create(s.key, nonce)
	s.encrypt = true
}

func (s *Stream) Send(data []byte, dataType Type) {
	if s.encrypt {
		log.Debug().Msg("Encrypting data")
		data = s.aesCipher.Encrypt(data)
	}
	header := Header{Len: uint16(len(data)), Type: dataType}
	n, err := s.Conn.Write(append(header.build(), data...))
	log.Debug().Int("len", n).Msg("Sending data to socket")
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Can't write to socket")
	}
}

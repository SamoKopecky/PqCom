package network

import (
	"encoding/binary"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/SamoKopecky/pqcom/main/io"
	"github.com/rs/zerolog/log"
)

type Type uint8

const (
	LEN_LEN       = 2
	TYPE_LEN      = 1
	HEADER_LEN    = LEN_LEN + TYPE_LEN
	KEM_TYPE_LEN  = 1
	SIGN_TYPE_LEN = 1
)

const (
	ClientInitT Type = iota
	ServerInitT
	ContentT
	ErrorT
)

type Header struct {
	Len  uint16
	Type Type
}

type ErrorMsg struct {
	reason string
}

type ClientInit struct {
	header   Header
	kemType  uint8
	signType uint8
	eK       []byte
	nonce    []byte
	sig      []byte
}

type ServerInit struct {
	header Header
	keyC   []byte
}

func (h *Header) parse(data []byte) {
	h.Len = binary.BigEndian.Uint16(cut(&data, LEN_LEN))
	h.Type = Type(data[0])
}

func (h *Header) build() []byte {
	headerLen := make([]byte, LEN_LEN)
	var headerType byte
	binary.BigEndian.PutUint16(headerLen, h.Len)
	headerType = byte(h.Type)
	return append(headerLen, headerType)
}

func (ci *ClientInit) parse(data []byte) []byte {
	log.Info().Msg("Parsing client init")
	var singedData []byte
	ci.kemType = cut(&data, KEM_TYPE_LEN)[0]
	ci.signType = cut(&data, SIGN_TYPE_LEN)[0]
	ci.eK = cut(&data, ekLen)
	ci.nonce = cut(&data, crypto.NONCE_LEN)
	ci.sig = data
	singedData = append(singedData, ci.kemType)
	singedData = append(singedData, ci.signType)
	singedData = append(singedData, ci.eK...)
	singedData = append(singedData, ci.nonce...)
	return singedData
}

func (ci *ClientInit) build() []byte {
	var data []byte
	data = append(data, ci.kemType)
	data = append(data, ci.signType)
	data = append(data, ci.eK...)
	data = append(data, ci.nonce...)
	if len(ci.sig) != 0 {
		data = append(data, ci.sig...)
	}
	return data
}

func (si *ServerInit) parse(data []byte) {
	log.Info().Msg("Parsing server init")
	si.keyC = data
}

func (si *ServerInit) build() []byte {
	return append([]byte{}, si.keyC...)
}

func (e *ErrorMsg) parse(data []byte) {
	log.Info().Msg("Parsing error msg")
	e.reason = string(data)
}

func (e *ErrorMsg) build() []byte {
	return []byte(e.reason)
}

func getNumberBytes(number, size int) []byte {
	data := make([]byte, size)
	binary.BigEndian.PutUint16(data, uint16(number))
	return data
}

func cut(data *[]byte, index int) []byte {
	if index > len(*data) {
		log.Error().
			Int("index", index).
			Int("data len", len(*data)).
			Msg("Error parsing header")
		return []byte{}
	}
	cut := (*data)[:index]
	*data = io.Copy((*data)[index:])
	return cut
}

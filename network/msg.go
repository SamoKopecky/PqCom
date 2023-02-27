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
	TIMESTAMP_LEN = 8
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
	header    Header
	kemType   uint8
	signType  uint8
	timestamp uint64
	eK        []byte
	nonce     []byte
	sig       []byte
}

type ServerInit struct {
	header Header
	keyC   []byte
}

func (h *Header) parse(data []byte) {
	h.Len = bytesToInt[uint16](2, (cut(&data, LEN_LEN)))
	h.Type = Type(data[0])
}

func (h *Header) build() []byte {
	headerLen := make([]byte, LEN_LEN)
	var headerType byte
	binary.BigEndian.PutUint16(headerLen, h.Len)
	headerType = byte(h.Type)
	return append(headerLen, headerType)
}

func (ci *ClientInit) parse(data []byte) {
	// TODO: Make more efficient
	log.Info().Msg("Parsing client init")

	ci.kemType = cut(&data, KEM_TYPE_LEN)[0]
	ci.signType = cut(&data, SIGN_TYPE_LEN)[0]
	timestampBytes := cut(&data, 8)
	ci.timestamp = bytesToInt[uint64](8, timestampBytes)
	ci.eK = cut(&data, ekLen)
	ci.nonce = cut(&data, crypto.NONCE_LEN)
	ci.sig = data
}

func (ci *ClientInit) build() []byte {
	data := ci.payload()
	if len(ci.sig) != 0 {
		data = append(data, ci.sig...)
	}
	return data
}

func (ci *ClientInit) payload() []byte {
	var payload []byte
	payload = append(payload, ci.kemType)
	payload = append(payload, ci.signType)
	payload = append(payload, intToBytes(int(ci.timestamp), 8)...)
	payload = append(payload, ci.eK...)
	payload = append(payload, ci.nonce...)
	return payload
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

func intToBytes(number, n int) []byte {
	data := make([]byte, n)
	if n == 2 {
		binary.BigEndian.PutUint16(data, uint16(number))
	} else if n == 4 {
		binary.BigEndian.PutUint32(data, uint32(number))
	} else if n == 8 {
		binary.BigEndian.PutUint64(data, uint64(number))
	}
	return data
}

func bytesToInt[T uint16 | uint32 | uint64](n int, data []byte) T {
	if n == 2 {
		return T(binary.BigEndian.Uint16(data))
	} else if n == 4 {
		return T(binary.BigEndian.Uint32(data))
	} else {
		return T(binary.BigEndian.Uint64(data))
	}
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

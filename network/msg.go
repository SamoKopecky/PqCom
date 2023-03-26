package network

import (
	"encoding/binary"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/SamoKopecky/pqcom/main/myio"
	"github.com/rs/zerolog/log"
)

type Type uint8

const (
	LenLen       = 2
	TypeLen      = 1
	HeaderLen    = LenLen + TypeLen
	KemTypeLen   = 1
	SignTypeLen  = 1
	TimestampLen = 8
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

type Error struct {
	reason string
}

type ClientInit struct {
	header    Header
	kemType   uint8
	signType  uint8
	timestamp uint64
	eK        []byte
	sig       []byte
}

type ServerInit struct {
	header Header
	keyC   []byte
	sig    []byte
}

type Content struct {
	header Header
	nonce  []byte
	data   []byte
}

func (h *Header) parse(data []byte) {
	h.Len = bytesToInt[uint16](2, (cut(&data, LenLen)))
	h.Type = Type(data[0])
}

func (h *Header) build() []byte {
	headerLen := make([]byte, LenLen)
	var headerType byte
	binary.BigEndian.PutUint16(headerLen, h.Len)
	headerType = byte(h.Type)
	return append(headerLen, headerType)
}

func (ci *ClientInit) parse(data []byte) {
	log.Info().Msg("Parsing client init")

	ci.kemType = cut(&data, KemTypeLen)[0]
	ci.signType = cut(&data, SignTypeLen)[0]
	timestampBytes := cut(&data, 8)
	ci.timestamp = bytesToInt[uint64](8, timestampBytes)
	ci.eK = cut(&data, ekLen)
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
	return payload
}

func (si *ServerInit) parse(data []byte) {
	log.Info().Msg("Parsing server init")
	si.keyC = cut(&data, cLen)
	si.sig = data
}

func (si *ServerInit) build() []byte {
	data := si.payload()
	if len(si.sig) != 0 {
		data = append(data, si.sig...)
	}
	return data
}

func (si *ServerInit) payload() []byte {
	return si.keyC
}

func (e *Error) parse(data []byte) {
	log.Info().Msg("Parsing error msg")
	e.reason = string(data)
}

func (e *Error) build() []byte {
	return []byte(e.reason)
}

func (c *Content) parse(data []byte) {
	log.Info().Msg("Parsing content")
	c.nonce = data[:crypto.NONCE_LEN]
	c.data = data[crypto.NONCE_LEN:]
}

func (c *Content) build() []byte {
	var data []byte
	data = append(data, c.nonce...)
	data = append(data, c.data...)
	return data
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
		log.Fatal().
			Int("index", index).
			Int("data_len", len(*data)).
			Msg("Error parsing header")
		return []byte{}
	}
	cut := (*data)[:index]
	*data = myio.Copy((*data)[index:])
	return cut
}

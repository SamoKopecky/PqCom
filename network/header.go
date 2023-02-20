package network

import (
	"encoding/binary"
)

type Type uint8

const (
	LEN_LEN    = 2
	TYPE_LEN   = 1
	HEADER_LEN = LEN_LEN + TYPE_LEN
	PK_LEN_LEN = 2
)

const (
	ClientInitT Type = iota
	ServerInitT
	ContentT
)

type Header struct {
	Len  uint16
	Type Type
}

type ClientInit struct {
	header Header
	pkLen  uint16
	pk     []byte
	nonce  []byte
}

type ServerInit struct {
	header Header
	keyC   []byte
}

func (header *Header) parse(data []byte) {
	header.Len = binary.BigEndian.Uint16(data[:LEN_LEN])
	header.Type = Type(data[LEN_LEN+TYPE_LEN-1])
}

func (clientInit *ClientInit) parse(data []byte) {
	clientInit.pkLen = binary.BigEndian.Uint16(data[:PK_LEN_LEN])
	clientInit.pk = data[PK_LEN_LEN : clientInit.pkLen+PK_LEN_LEN]
	clientInit.nonce = data[clientInit.pkLen+PK_LEN_LEN:]
}

func (serverInit *ServerInit) parse(data []byte) {
	serverInit.keyC = data
}

func (header *Header) build() []byte {
	headerLen := make([]byte, LEN_LEN)
	var headerType byte
	binary.BigEndian.PutUint16(headerLen, header.Len)
	headerType = byte(header.Type)
	return append(headerLen, headerType)
}

func (clientInit *ClientInit) build() []byte {
	data := make([]byte, PK_LEN_LEN)
	binary.BigEndian.PutUint16(data, clientInit.pkLen)
	data = append(data, clientInit.pk...)
	data = append(data, clientInit.nonce...)
	return data
}

func (serverInit *ServerInit) build() []byte {
	return append([]byte{}, serverInit.keyC...)
}

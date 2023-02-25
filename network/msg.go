package network

import (
	"encoding/binary"

	"github.com/SamoKopecky/pqcom/main/crypto"
)

type Type uint8

const (
	LEN_LEN     = 2
	TYPE_LEN    = 1
	HEADER_LEN  = LEN_LEN + TYPE_LEN
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
	eK     []byte
	nonce  []byte
	sig    []byte
}

type ServerInit struct {
	header Header
	keyC   []byte
}

func (header *Header) parse(data []byte) {
	header.Len = binary.BigEndian.Uint16(data[:LEN_LEN])
	header.Type = Type(data[LEN_LEN+TYPE_LEN-1])
}

func (clientInit *ClientInit) parse(data []byte) []byte {
	var singed_data []byte
	clientInit.eK = cut(&data, ekLen)
	clientInit.nonce = cut(&data, crypto.NONCE_LEN)
	clientInit.sig = data
	singed_data = append(singed_data, clientInit.eK...)
	singed_data = append(singed_data, clientInit.nonce...)
	return singed_data
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
	var data []byte
	data = append(data, clientInit.eK...)
	data = append(data, clientInit.nonce...)
	if len(clientInit.sig) != 0 {
		data = append(data, clientInit.sig...)
	}
	return data
}

func (serverInit *ServerInit) build() []byte {
	return append([]byte{}, serverInit.keyC...)
}

func getNumberBytes(number, size int) []byte {
	data := make([]byte, size)
	binary.BigEndian.PutUint16(data, uint16(number))
	return data
}

func cut(data *[]byte, index int) []byte {
	cut := (*data)[:index]
	*data = (*data)[index:]
	return cut
}

package network

import (
	"encoding/binary"
)

const LEN_BYTES = 2
const HEADER_LEN = LEN_BYTES

type Header struct {
	Len uint16 // 2 B
}

func (header *Header) parse(data []byte) {
	header.Len = binary.BigEndian.Uint16(data[:LEN_BYTES])
}

func (header *Header) build() []byte {
	data := make([]byte, HEADER_LEN)
	binary.BigEndian.PutUint16(data, header.Len)
	return data
}

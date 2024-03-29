package myio

import (
	"bufio"
	"io"
	"math/rand"
	"reflect"

	"github.com/rs/zerolog/log"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ReadByChunks(reader io.Reader, chunks chan<- []byte, chunkSize int) {
	r := bufio.NewReader(reader)
	buf := make([]byte, 0, chunkSize)
	for {
		n, err := r.Read(buf[:cap(buf)])
		log.Debug().
			Int("len", n).
			Str("reader", reflect.TypeOf(reader).String()).
			Msg("Reading data")
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Error().Str("error", err.Error()).Msg("Error while reading 0 bytes")
		}
		if err != nil && err != io.EOF {
			log.Error().Str("error", err.Error()).Msg("Error reading")
		}
		chunks <- Copy(buf[:n])
	}
}

func Read(r io.Reader, buf []byte) (n int, err error) {
	n, err = r.Read(buf[:cap(buf)])
	log.Debug().Int("len", n).Msg("Reading data")
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Error reading")
	}
	return
}

func RandStringBytes(n int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func Copy(src []byte) []byte {
	cpy := make([]byte, len(src))
	copy(cpy, src)
	return cpy
}

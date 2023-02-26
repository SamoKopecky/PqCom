package io

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/rs/zerolog/log"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ReadByChunks(reader io.Reader, chunks chan<- []byte, chunkSize int) {
	r := bufio.NewReader(reader)
	buf := make([]byte, 0, chunkSize)
	for {
		// TODO: Make is that buf doesn't have to initialize every time
		n, err := r.Read(buf[:cap(buf)])
		log.Debug().
			Int("len", n).
			Str("reader", reflect.TypeOf(reader).Name()).
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
		log.Debug().Msg("sending to channel")
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

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ContainsDir(file string, dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("Error reading directory")
	}
	for _, e := range entries {
		if file == e.Name() {
			return true
		}
	}
	return false
}

func ReadUserInput(promt string) string {
	fmt.Print(promt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func Copy(src []byte) []byte {
	cpy := make([]byte, len(src))
	copy(cpy, src)
	return cpy
}

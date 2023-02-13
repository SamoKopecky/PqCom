package handler

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ReadByChunks(reader io.Reader, chunks chan<- []byte, chunkSize int) {
	r := bufio.NewReader(reader)
	for {
		// TODO: Make is that buf doesn't have to initialize every time
		buf := make([]byte, 0, chunkSize)
		n, err := r.Read(buf[:cap(buf)])
		log.WithFields(log.Fields{
			"len":    n,
			"reader": reflect.TypeOf(reader),
		}).Debug("Reading data")
		chunks <- buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.WithField("error", err).Error("Error while reading 0 bytes")
		}
		if err != nil && err != io.EOF {
			log.WithField("error", err).Error("Error reading")
		}
	}
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func ContainsDir(file string, dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.WithField("error", err).Error("Error reading directory")
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

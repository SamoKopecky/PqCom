package app

import (
	"fmt"
	"io"
	"os"

	"github.com/SamoKopecky/pqcom/main/myio"
	"github.com/SamoKopecky/pqcom/main/network"
	"github.com/rs/zerolog/log"
)

func Chat(destAddr string, srcPort, destPort int, connect bool, logFile *os.File) {
	if connect {
		stream := network.Connect(destAddr, destPort)
		ChatTui(stream, stream.Send)
	} else {
		streamFactory := make(chan network.Stream)
		go network.Listen(srcPort, streamFactory, false)
		stream := <-streamFactory
		ChatTui(stream, stream.Send)
	}
	fmt.Printf("Logs saved to a file at %s\n", logFile.Name())
}

func Send(destAddr string, srcPort, destPort int, filePath string) {
	stream := network.Connect(destAddr, destPort)
	chunks := make(chan []byte)
	var source io.Reader
	var err error

	if filePath != "" {
		source, err = os.Open(filePath)
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("Error opening file")
		}
		log.Debug().Msg("Using file as data source")
	} else {
		source = os.Stdin
		log.Debug().Msg("Using stdin as data source")
	}
	go func() {
		myio.ReadByChunks(source, chunks, network.CHUNK_SIZE)
		close(chunks)
	}()
	for msg := range chunks {
		stream.Send(msg, network.ContentT)
	}
	log.Info().Str("addr", stream.Conn.RemoteAddr().String()).Msg("Done sending")
}

func Receive(destAddr string, srcPort, destPort int, dir string) {
	streamFactory := make(chan network.Stream)
	go network.Listen(srcPort, streamFactory, true)
	for {
		stream := <-streamFactory
		if dir != "" {
			go dirFileWriter(stream.Msg, dir)
		} else {
			go printer(stream)
		}
	}
}

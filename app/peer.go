package app

import (
	"io"
	"os"

	myio "github.com/SamoKopecky/pqcom/main/io"
	"github.com/SamoKopecky/pqcom/main/network"
	log "github.com/sirupsen/logrus"
)

func Chat(destAddr string, srcPort, destPort int, connect bool) {
	if connect {
		stream := network.Connect(destAddr, destPort)
		go printer(stream, false)
		for {
			data := []byte(myio.ReadUserInput(""))
			stream.Send(data)
		}
	} else {
		streamFactory := make(chan network.Stream)
		go network.Listen(srcPort, streamFactory, false)
		stream := <-streamFactory
		go printer(stream, false)
		for {
			data := []byte(myio.ReadUserInput(""))
			stream.Send(data)
		}
	}
}

func Send(destAddr string, srcPort, destPort int, filePath string) {
	stream := network.Connect(destAddr, destPort)
	chunks := make(chan []byte)
	var source io.Reader
	var err error

	if filePath != "" {
		source, err = os.Open(filePath)
		if err != nil {
			log.WithField("error", err).Error("Error opening file")
		}
	} else {
		source = os.Stdin
	}
	go func() {
		myio.ReadByChunks(source, chunks, network.ChunkSize)
		close(chunks)
	}()
	for msg := range chunks {
		stream.Send(msg)
	}
	log.WithField("addr", stream.Conn.RemoteAddr()).Info("Done sending")
}

func Receive(destAddr string, srcPort, destPort int, dir string) {
	streamFactory := make(chan network.Stream)
	go network.Listen(srcPort, streamFactory, true)
	for {
		stream := <-streamFactory
		if dir != "" {
			go dirFileWriter(stream.Data, dir)
		} else {
			go printer(stream, true)
		}
	}
}

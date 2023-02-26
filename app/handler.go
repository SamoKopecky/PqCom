package app

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SamoKopecky/pqcom/main/io"
	"github.com/SamoKopecky/pqcom/main/network"
	"github.com/rs/zerolog/log"
)

const sufix = 5

func dirFileWriter(recv <-chan network.Msg, dir string) {
	newFile := true
	var fileName string
	var file *os.File
	var err error

	for {
		msg := <-recv

		if newFile {
			for io.ContainsDir(fileName, dir) || fileName == "" {
				fileName = fmt.Sprint("pqcom_temp_", io.RandStringBytes(sufix))
			}
			filePath := fmt.Sprint(dir, string(os.PathSeparator), fileName)

			file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Error().Str("error", err.Error()).Msg("Error opening file")
			}
			newFile = false
		}

		w := bufio.NewWriter(file)
		n, err := w.Write(msg.Data)
		log.Debug().Int("len", n).Msg("Write data to file")
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("Error writing to file")
		}
		if n == 0 {
			newFile = true
			err = file.Close()
			if err != nil {
				log.Error().Str("error", err.Error()).Msg("Error closing file")
			}
		}
		w.Flush()
	}
}

func printer(stream network.Stream, clean bool) {
	var msg network.Msg
	for {
		msg = <-stream.Msg
		if clean {
			fmt.Printf("%s", string(msg.Data))
			continue
		}
		fmt.Printf("[%s]: %s", stream.Conn.RemoteAddr(), string(msg.Data))
	}
}

package handler

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const sufix = 5

func FileWriter(dir string, recv <-chan []byte) {
	newFile := true
	var fileName string
	var file *os.File
	var err error

	for {
		msg := <-recv

		if newFile {
			for ContainsDir(fileName, dir) || fileName == "" {
				fileName = fmt.Sprint("pqcom_temp_", randStringBytes(sufix))
			}
			filePath := fmt.Sprint(dir, string(os.PathSeparator), fileName)

			file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.WithField("error", err).Error("Error opening file")
			}

			newFile = false
		}

		w := bufio.NewWriter(file)
		n, err := w.Write(msg)
		log.WithField("len", n).Debug("Write data to file")
		if err != nil {
			log.WithField("error", err).Error("Error writing to file")
		}
		if n == 0 {
			newFile = true
			err = file.Close()
			if err != nil {
				log.WithField("error", err).Error("Error closing file")
			}
		}
		w.Flush()
	}
}

func Printer(recv <-chan []byte, clean bool) {
	for {
		msg := <-recv
		if clean {
			fmt.Printf("%s", string(msg))
			continue
		}
		fmt.Printf("[%s]: %s", "temp", string(msg))
	}
}

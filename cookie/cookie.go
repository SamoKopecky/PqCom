package cookie

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/SamoKopecky/pqcom/main/io"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/sha3"
)

type Cookie struct {
	Seed      []byte
	Timestamp uint64
}

func Get() uint64 {
	return uint64(time.Now().UnixMicro())
}

func (c *Cookie) Save() {
	name, location := c.fileNameAndDir()
	filePath := fmt.Sprintf("%s%s", location, name)

	_, err := os.Stat(location)
	if err != nil {
		log.Info().Str("dir", filePath).Msg("Creating directory")
		os.Mkdir(location, 0700)
	}
	flags := os.O_WRONLY
	if !c.Exists() {
		log.Info().Str("name", name).Msg("Crating new cookie")
		flags = flags | os.O_CREATE
	}
	file, err := os.OpenFile(filePath, flags, 0600)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening file for writing")
	}
	log.Info().Str("name", name).Msg("Writing timestamp to cookie file")
	toWrite := make([]byte, 8)
	binary.BigEndian.PutUint64(toWrite, c.Timestamp)
	file.Write(toWrite)
	file.Close()
}

func (c *Cookie) IsNewer() bool {
	name, location := c.fileNameAndDir()
	filePath := fmt.Sprintf("%s%s", location, name)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening file for reading")
	}
	buf := make([]byte, 8)
	_, err = file.Read(buf)
	oldTimestamp := binary.BigEndian.Uint64(buf)
	return oldTimestamp <= c.Timestamp
}

func (c *Cookie) Exists() bool {
	name, location := c.fileNameAndDir()
	contains, err := io.ContainsFile(name, location)
	_, is := err.(*os.PathError)
	if is {
		return false
	}
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening directory")
	}
	return contains
}

func (c *Cookie) fileNameAndDir() (name, dir string) {
	dir = io.HomeSubDir(".cache")
	hash := sha3.Sum512(c.Seed)
	name = io.RandStringBytes(16, int64(binary.BigEndian.Uint64(hash[:])))
	return
}

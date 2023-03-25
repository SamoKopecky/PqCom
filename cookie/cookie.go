package cookie

import (
	"encoding/binary"
	"os"
	"time"

	"github.com/SamoKopecky/pqcom/main/myio"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/sha3"
)

const cookieSize = 8

type Cookie struct {
	Seed      []byte
	Timestamp uint64
}

func Get() uint64 {
	return uint64(time.Now().UnixMicro())
}

func (c *Cookie) Save() {
	name, dir := c.nameAndDir()
	filePath := dir + name

	myio.CreatePath(filePath)

	flags := os.O_WRONLY
	if !c.Exists() {
		log.Info().Str("name", name).Msg("Crating new cookie")
		flags = flags | os.O_CREATE
	}
	file, err := os.OpenFile(filePath, flags, myio.NewFilePerms)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening file for writing")
	}
	log.Info().Str("name", name).Msg("Writing timestamp to cookie file")
	toWrite := make([]byte, cookieSize)
	binary.BigEndian.PutUint64(toWrite, c.Timestamp)
	file.Write(toWrite)
	file.Close()
}

func (c *Cookie) IsNewer() bool {
	name, dir := c.nameAndDir()
	file, err := os.Open(dir + name)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening cookie for reading")
	}
	buf := make([]byte, cookieSize)
	_, err = file.Read(buf)
	oldTimestamp := binary.BigEndian.Uint64(buf)
	return oldTimestamp <= c.Timestamp
}

func (c *Cookie) Exists() bool {
	contains, err := myio.ContainsFile(c.nameAndDir())
	_, is := err.(*os.PathError)
	if is {
		return false
	}
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Error opening directory")
	}
	return contains
}

func (c *Cookie) nameAndDir() (name, dir string) {
	dir = myio.HomeSubDir(myio.Cookie)
	hash := sha3.Sum512(c.Seed)
	name = myio.RandStringBytes(16, int64(binary.BigEndian.Uint64(hash[:])))
	return
}

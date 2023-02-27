package timestamp

import (
	"time"
)

func Get() uint64 {
	return uint64(time.Now().UnixMicro())
}

func Save(timestamp uint64) {
	// TODO:
}

func Check(timestamp uint64) {
	// TODO:
}

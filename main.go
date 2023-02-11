package main

import (
	"io"

	"github.com/SamoKopecky/pqcom/main/benchmark"
	"github.com/SamoKopecky/pqcom/main/network"
	log "github.com/sirupsen/logrus"
)

func main() {
	o := options{}
	o.parseArgs()

	if !o.log {
		log.SetOutput(io.Discard)
	}
	if o.benchmark {
		benchmark.Run(o.iterations)
	}
	if o.app {
		network.Start(o.destAddr, o.srcPort, o.destPort, o.stdin)
	}
}

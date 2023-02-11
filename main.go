package main

import (
	"flag"
	"io"

	"github.com/SamoKopecky/pqcom/main/benchmark"
	"github.com/SamoKopecky/pqcom/main/network"
	log "github.com/sirupsen/logrus"
)

type options struct {
	destPort   int
	srcPort    int
	iterations int
	destAddr   string
	stdin      bool
	benchmark  bool
	app        bool
	log        bool
}

func (o *options) parseArgs() {
	flag.BoolVar(&o.benchmark, "b", false, "Bbenchmark")
	flag.BoolVar(&o.log, "l", false, "Enable logging")
	flag.BoolVar(&o.app, "a", false, "Start app")
	flag.BoolVar(&o.stdin, "si", false, "Read from stdin")

	flag.IntVar(&o.iterations, "i", 1000, "Number of iterations")
	flag.IntVar(&o.srcPort, "sp", 4040, "Source port")
	flag.IntVar(&o.destPort, "dp", 4040, "Destination port")

	flag.StringVar(&o.destAddr, "da", "localhost", "Destination address")

	flag.Parse()
}

func main() {
	var o options
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

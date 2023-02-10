package main

import (
	"flag"
	"io"

	"github.com/SamoKopecky/pqcom/main/benchmark"
	"github.com/SamoKopecky/pqcom/main/network"
	log "github.com/sirupsen/logrus"
)

type options struct {
	dport      int
	lport      int
	iterations int
	daddr      string
	stdin      bool
	benchmark  bool
	app        bool
	log        bool
}

func (o *options) parseArgs() {
	flag.BoolVar(&o.benchmark, "b", false, "srart benchmark")
	flag.BoolVar(&o.log, "log", false, "enable logging")
	flag.BoolVar(&o.app, "a", false, "start app")
	flag.BoolVar(&o.stdin, "si", false, "read from stdin")

	flag.IntVar(&o.iterations, "i", 1000, "set the number of iterations")
	flag.IntVar(&o.lport, "lp", 4040, "local port")
	flag.IntVar(&o.dport, "dp", 4040, "destination port")

	flag.StringVar(&o.daddr, "da", "localhost", "destination address")

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
		network.Start(o.daddr, o.lport, o.dport, o.stdin)
	}
}

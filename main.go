package main

import (
	"flag"

	"github.com/SamoKopecky/pqcom/main/benchmark"
	"github.com/SamoKopecky/pqcom/main/network"
)

type options struct {
	dport, lport   int
	daddr, laddr   string
	benchmark, app bool
	iterations     int
}

func (o *options) parseArgs() {
	flag.BoolVar(&o.benchmark, "b", false, "srart benchmark")
	flag.IntVar(&o.iterations, "i", 1000, "set the number of iterations")

	flag.BoolVar(&o.app, "a", false, "start app")
	flag.StringVar(&o.laddr, "la", "localhost", "local address")
	flag.StringVar(&o.daddr, "da", "localhost", "destination address")
	flag.IntVar(&o.lport, "lp", 4040, "local port")
	flag.IntVar(&o.dport, "dp", 4040, "destination port")
	flag.Parse()
}

func main() {
	var o options
	o.parseArgs()

	if o.benchmark {
		benchmark.Run(o.iterations)
	}
	if o.app {
		network.Start(o.laddr, o.daddr, o.lport, o.dport)
	}
}

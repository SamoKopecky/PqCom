package main

import "flag"

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
	flag.IntVar(&o.iterations, "i", 1000, "Number of iterations")
	flag.IntVar(&o.srcPort, "sp", 4040, "Source port")
	flag.IntVar(&o.destPort, "dp", 4040, "Destination port")

	flag.StringVar(&o.destAddr, "da", "localhost", "Destination address")

	flag.BoolVar(&o.benchmark, "b", false, "Bbenchmark")
	flag.BoolVar(&o.log, "l", false, "Enable logging")
	flag.BoolVar(&o.app, "a", false, "Start app")
	flag.BoolVar(&o.stdin, "si", false, "Read from stdin")

	flag.Parse()
}

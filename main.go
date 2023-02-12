package main

import "github.com/SamoKopecky/pqcom/main/cmd"

// func main() {
// 	o := options.Options{}
// 	o.ParseArgs()

// 	if !o.Log {
// 		log.SetOutput(io.Discard)
// 	} else if o.Benchmark {
// 		benchmark.Run(o.Iterations)
// 	} else {
// 		// network.Start(o.DestAddr, o.SrcPort, o.DestPort, o.Stdin, o.File, o.Chat, o.FilePath)
// 	}

// }

func main() {
	cmd.Execute()
}

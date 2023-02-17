package main

import (
	"github.com/SamoKopecky/pqcom/main/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	cmd.Execute()
}

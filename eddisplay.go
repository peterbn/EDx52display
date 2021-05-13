package main

import (
	"fmt"
	"log"

	"github.com/peterbn/EDx52display/conf"
	"github.com/peterbn/EDx52display/edreader"
	"github.com/peterbn/EDx52display/edsm"
	"github.com/peterbn/EDx52display/mfd"
)

func main() {
	conf := conf.LoadConf()

	err := mfd.InitDevice(edreader.DisplayPages, edsm.ClearCache)
	if err != nil {
		panic(err)
	}
	defer mfd.DeInitDevice()

	edreader.Start(conf)
	defer edreader.Stop()

	log.Println("EDx52Display running. Press enter to close.")
	fmt.Scanln() // keep it running until I get input
}

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/peterbn/EDx52display/conf"
	"github.com/peterbn/EDx52display/edreader"
	"github.com/peterbn/EDx52display/mfd"
)

func main() {
	conf := conf.LoadConf()

	edreader.Start(conf)
	defer edreader.Stop()

	// Start the file monitor for updates
	cmd := exec.Command("X52MFDDriver.exe", ".\\"+mfd.Filename)
	cmd.Stdout = os.Stdout // ensure the driver's output is also sent to the console
	cmd.Start()

	fmt.Scanln() // keep it running until I get input
}

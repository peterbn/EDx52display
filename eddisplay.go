package main

import (
	"fmt"
	"os/exec"

	"github.com/peterbn/EDx52display/conf"
	"github.com/peterbn/EDx52display/edreader"
	"github.com/peterbn/EDx52display/mfd"
)

func main() {
	conf := conf.LoadConf()
	edreader.Start(conf)
	defer edreader.Stop()

	var mfdData = mfd.Display{
		Pages: []mfd.Page{mfd.Page{Lines: []string{"Time since start", "0s"}}},
	}
	mfd.Write(mfdData)

	// Start the file monitor for updates
	cmd := exec.Command("X52MFDDriver.exe", ".\\"+mfd.Filename)
	cmd.Start()
	defer cmd.Process.Kill()

	fmt.Scanln() // keep it running until I get input
}

package main

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"

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
	cmd.Start()
	log.Println("EDx52Display running. Press enter to close.")
	fmt.Scanln() // keep it running until I get input
	cmd.Process.Signal(syscall.SIGTERM)
	cmd.Process.Kill() // and then kill the child process to clean up.
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"github.com/peterbn/EDx52display/conf"
	"github.com/peterbn/EDx52display/edreader"
	"github.com/peterbn/EDx52display/edsm"
	"github.com/peterbn/EDx52display/mfd"
)

func main() {
	conf := conf.LoadConf()

	edreader.Start(conf)
	defer edreader.Stop()

	// Start the file monitor for updates
	cmd := exec.Command("X52MFDDriver.exe", ".\\"+mfd.Filename)

	messages, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		reader := bufio.NewReader(messages)
		for {
			line, _ := reader.ReadString('\n')
			if strings.HasPrefix(line, ">MFD-SELECT PRESSED<") {
				log.Println("Got request to update EDSM cache")
				// Reset the EDSM cache and reload the display
				edsm.ClearCache()
			}
		}
	}()

	cmd.Start()
	log.Println("EDx52Display running. Press enter to close.")
	fmt.Scanln() // keep it running until I get input
	cmd.Process.Signal(syscall.SIGTERM)
	cmd.Process.Kill() // and then kill the child process to clean up.
}

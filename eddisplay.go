package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {

	var mfd = MfdDisplay{
		Pages: []MfdPage{MfdPage{Lines: []string{"Time since start", "0s"}}},
	}
	writeMFD(mfd)

	const delay = 1
	ticker := time.NewTicker(delay * time.Second)
	go func() {

		var i = 0

		for range ticker.C {
			i = i + delay
			dString := fmt.Sprintf("%ds", i)
			var mfdInner = MfdDisplay{
				Pages: []MfdPage{MfdPage{Lines: []string{"Time since start", dString}}},
			}

			writeMFD(mfdInner)
		}
	}()
	// Start the file monitor for updates
	cmd := exec.Command("X52MFDDriver.exe", ".\\"+mfdFilename)
	cmd.Start()
	defer cmd.Process.Kill()

	fmt.Scanln() // keep it running until I get input
}

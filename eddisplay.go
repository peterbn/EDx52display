package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const filename = "mfd.json"

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
	cmd := exec.Command("X52MFDDriver.exe", ".\\"+filename)
	cmd.Start()
	defer cmd.Process.Kill()

	fmt.Scanln() // keep it running until I get input
}

func writeMFD(mfd MfdDisplay) {
	data, err := json.Marshal(mfd)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(filename)
	if err != nil {

		fmt.Println(err)
		return

	}
	defer f.Close()

	f.Write(data)
	f.Sync()
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const mfdFilename = "mfd.json"

// MfdDisplay is the main display structure to write
type MfdDisplay struct {
	Pages []MfdPage `json:"pages"`
}

// MfdPage is a single page of information to show on the MFD
type MfdPage struct {
	Lines []string `json:"lines"`
}

func writeMFD(mfd MfdDisplay) {
	data, err := json.Marshal(mfd)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(mfdFilename)
	if err != nil {

		fmt.Println(err)
		return

	}
	defer f.Close()

	f.Write(data)
	f.Sync()
}

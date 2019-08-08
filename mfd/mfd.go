package mfd

import (
	"encoding/json"
	"fmt"
	"os"
)

// Filename is the name of the file written with MFD data
const Filename = "mfd.json"

// Display is the main display structure to write
type Display struct {
	Pages []Page `json:"pages"`
}

// Page is a single page of information to show on the MFD
type Page struct {
	Lines []string `json:"lines"`
}

// Write writes the MFD file
func Write(mfd Display) {
	data, err := json.Marshal(mfd)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(Filename)
	if err != nil {

		fmt.Println(err)
		return

	}
	defer f.Close()

	f.Write(data)
	f.Sync()
}

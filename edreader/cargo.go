package edreader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/peterbn/EDx52display/mfd"
)

// FileCargo is the name of the processed Cargo file
const FileCargo = "Cargo.json"

type Cargo struct {
	Count     int
	Inventory []CargoLine
}

type CargoLine struct {
	Name          string
	Count         int
	Stolen        int
	NameLocalized string `json:"Name_Localised"`
}

func handleCargoFile(file string) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	var cargo Cargo
	json.Unmarshal(data, &cargo)

	display := []string{}
	display = append(display, fmt.Sprintf("Cargo: %d", cargo.Count))
	for _, line := range cargo.Inventory {
		name := line.Name
		if line.NameLocalized != "" {
			name = line.NameLocalized
		}
		cline := fmt.Sprintf("%s: %d", name, line.Count)
		display = append(display, cline)
	}

	var mfdData = mfd.Display{
		Pages: []mfd.Page{mfd.Page{Lines: display}},
	}

	mfd.Write(mfdData)
}

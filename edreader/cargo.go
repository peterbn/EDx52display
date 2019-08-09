package edreader

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/peterbn/EDx52display/mfd"
)

// FileCargo is the name of the processed Cargo file
const FileCargo = "Cargo.json"

const (
	nameFileFolder        = "./names/"
	commodityNameFile     = nameFileFolder + "commodity.csv"
	rareCommodityNameFile = nameFileFolder + "rare_commodity.csv"
)

// Cargo struct to load the cargo file saved by ED
type Cargo struct {
	Count     int
	Inventory []CargoLine
}

// CargoLine struct to load individual items in the cargo
type CargoLine struct {
	Name          string
	Count         int
	Stolen        int
	NameLocalized string `json:"Name_Localised"`
}

var names map[string]string

func init() {
	initNameMap()
}

func handleCargoFile(file string) {
	log.Println("Loading cargo from file", file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	var cargo Cargo
	json.Unmarshal(data, &cargo)

	display := []string{}
	display = append(display, fmt.Sprintf("Cargo: %d", cargo.Count))
	display = append(display, renderCargo(cargo)...)

	Mfd.Pages[pageCargo].Lines = display

	mfd.Write(Mfd)
}

func renderCargo(cargo Cargo) []string {
	display := []string{}

	for _, line := range cargo.Inventory {
		name := line.Name
		displayName, ok := names[strings.ToLower(name)]
		if ok {
			name = displayName
		}
		cline := fmt.Sprintf("%s: %d", name, line.Count)
		display = append(display, cline)
	}

	sort.Strings(display)

	return display
}

func initNameMap() {

	commodity := readCsvFile(commodityNameFile)
	rareCommodity := readCsvFile(rareCommodityNameFile)

	names = make(map[string]string)

	mapCommodities := func(comms [][]string) {
		for _, com := range comms[1:] { //skipping the header line
			symbol := com[1]
			symbol = strings.ToLower(symbol)
			name := com[3]
			names[symbol] = name
		}
	}
	mapCommodities(commodity)
	mapCommodities(rareCommodity)
}

func readCsvFile(filename string) [][]string {
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Panicln(err)
		return nil
	}
	defer csvfile.Close()
	csvreader := csv.NewReader(csvfile)
	records, err := csvreader.ReadAll()
	if err != nil {
		log.Panicln(err)
		return nil
	}
	return records
}

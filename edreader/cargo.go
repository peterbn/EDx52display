package edreader

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

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

func (cl CargoLine) displayname() string {
	name := cl.Name
	displayName, ok := names[strings.ToLower(name)]
	if ok {
		name = displayName
	}
	return name
}

var names map[string]string

func init() {
	log.Debugln("Initializing cargo name map...")
	initNameMap()
}

func handleCargoFile(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	var cargo Cargo
	json.Unmarshal(data, &cargo)

	page := mfd.NewPage()

	renderCargo(&page, cargo)

	MfdLock.Lock()
	Mfd.Pages[pageCargo] = page
	MfdLock.Unlock()
}

func renderCargo(page *mfd.Page, cargo Cargo) {
	page.Add("#Cargo: %03d/%03d#", cargo.Count, ModulesInfoCargoCapacity())
	sort.Slice(cargo.Inventory, func(i, j int) bool {
		a := cargo.Inventory[i]
		b := cargo.Inventory[j]
		return a.displayname() < b.displayname()
	})

	for _, line := range cargo.Inventory {
		page.Add("%s: %d", line.displayname(), line.Count)
	}
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

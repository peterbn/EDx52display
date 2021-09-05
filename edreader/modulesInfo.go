package edreader

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

const FileModulesInfo = "ModulesInfo.json"

// ModulesInfo struct to load the ModulesInfo file saved by ED
type ModulesInfo struct {
	Modules []ModulesLine
}

// ModulesLine struct to load individual module in the ModuleInfo
type ModulesLine struct {
	Slot        string
	Item        string
}

var currentModules ModulesInfo

func handleModulesInfoFile(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errorln(err)
		return
	}

	json.Unmarshal(data, &currentModules)
}

func ModulesInfoCargoCapacity() int {
	cargoCapacity := 0

	for _, line := range currentModules.Modules {
		switch line.Item {
			case "int_cargorack_size1_class1":
				cargoCapacity += 2
				break
			case "int_cargorack_size2_class1":
				cargoCapacity += 4
				break
			case "int_cargorack_size3_class1":
				cargoCapacity += 8
				break
			case "int_cargorack_size4_class1":
				cargoCapacity += 16
				break
			case "int_cargorack_size5_class1":
				cargoCapacity += 32
				break
			case "int_cargorack_size6_class1":
				cargoCapacity += 64
				break
			case "int_cargorack_size7_class1":
				cargoCapacity += 128
				break
			case "int_cargorack_size8_class1":
				cargoCapacity += 256
				break
		}
	}

	return cargoCapacity
}

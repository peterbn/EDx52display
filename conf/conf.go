package conf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/sys/windows/registry"

	"gopkg.in/yaml.v2"
)

// Conf is the app config
type Conf struct {
	JournalsFolder    string
	RefreshRateMS     int
	EDSMCommanderName string
	EDSMAPIKey        string
}

// LoadConf loads the config from the yaml file
func LoadConf() Conf {
	data, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var conf Conf
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalln(err)
	}

	return conf
}

// ExpandJournalFolderPath expands any env variables in the journal folder path.
func (c Conf) ExpandJournalFolderPath() string {
	exp, _ := registry.ExpandString(c.JournalsFolder)
	return exp
}

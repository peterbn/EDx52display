package main

import (
	"flag"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/peterbn/EDx52display/conf"
	"github.com/peterbn/EDx52display/edreader"
	"github.com/peterbn/EDx52display/edsm"
	"github.com/peterbn/EDx52display/mfd"
)

// TextLogFormatter gives me custom command-line formatting
type TextLogFormatter struct{}

func (f *TextLogFormatter) Format(entry *log.Entry) ([]byte, error) {
	level := entry.Level.String()
	message := entry.Message

	return []byte(strings.ToUpper(level) + ": " + message + "\n"), nil
}

func main() {

	defer func() {
		// Attempt to catch any crash messages before the cmd window closes
		if r := recover(); r != nil {
			log.Warnln("Crashed with message")
			log.Warnln(r)
			log.Warnln("Press RETURN to exit")
			fmt.Scanln() // keep it running until I get input
		}
	}()
	var logLevelArg string
	flag.StringVar(&logLevelArg, "log", "info", "Desired log level. One of [panic, fatal, error, warning, info, debug, trace]. Default: info.")

	flag.Parse()
	logLevel, err := log.ParseLevel(logLevelArg)
	if err != nil {
		log.Panic(err)
	}

	log.SetLevel(logLevel)

	log.SetFormatter(&TextLogFormatter{})

	conf := conf.LoadConf()

	err = mfd.InitDevice(edreader.DisplayPages, edsm.ClearCache)
	if err != nil {
		log.Panic(err)
	}
	defer mfd.DeInitDevice()

	edreader.Start(conf)
	defer edreader.Stop()

	log.Infoln("EDx52Display running. Press enter to close.")
	fmt.Scanln() // keep it running until I get input
}

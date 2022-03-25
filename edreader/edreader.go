package edreader

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbn/EDx52display/mfd"

	"github.com/peterbn/EDx52display/conf"
)

const DisplayPages = 3

var tick time.Ticker

const (
	pageCargo = iota
	pageLocation
	pageTargetInfo
)

// Mfd is the MFD display structure that will be used by this module. The number of pages should not be changed
var Mfd = mfd.Display{Pages: make([]mfd.Page, DisplayPages)}

// MfdLock locks the current MFD for reads and writes
var MfdLock = sync.RWMutex{}

// PrevMfd is the previous Mfd written to file, to be used for comparisons and avoid superflous updates.
var PrevMfd = Mfd.Copy()

// Start starts the Elite Dangerous journal reader routine
func Start(cfg conf.Conf) {
	// Update immediately, to ensure the mfd.json file exist
	log.Info("Starting journal listener")
	journalfolder := cfg.ExpandJournalFolderPath()
	log.Debugln("Looking for journal files in " + journalfolder)
	updateMFD(journalfolder)
	tick := time.NewTicker(time.Duration(cfg.RefreshRateMS) * time.Millisecond)

	go func() {
		for range tick.C {
			updateMFD(journalfolder)
		}
	}()
}

func updateMFD(journalfolder string) {
	// Read in the files at start before we start watching, to initialize
	journalFile := findJournalFile(journalfolder)
	handleJournalFile(journalFile)

	handleModulesInfoFile(filepath.Join(journalfolder, FileModulesInfo))
	handleCargoFile(filepath.Join(journalfolder, FileCargo))
	swapMfd()
}

// Stop closes the watcher again
func Stop() {
	tick.Stop()
}

func findJournalFile(folder string) string {
	// Based on https://github.com/EDCD/EDMarketConnector/blob/693463d3a0dbe58a1f72b83fc09a44a4398af3fd/monitor.py#L264
	// because I don't have Odyssey myself
	files, _ := filepath.Glob(filepath.Join(folder, "Journal.*.*.log"))

	var lastFileTime time.Time
	var mostRecentJournal = ""

	for _, filename := range files {
		info, err := os.Stat(filename)
		if err != nil {
			continue
		}
		if mostRecentJournal == "" || info.ModTime().After(lastFileTime) {
			lastFileTime = info.ModTime()
			mostRecentJournal = filename
		}
	}
	return mostRecentJournal
}

func swapMfd() {
	MfdLock.RLock()
	defer MfdLock.RUnlock()
	eq := cmp.Equal(Mfd, PrevMfd)
	if !eq {
		mfd.Write(Mfd)
		PrevMfd = Mfd.Copy()
	}
}

package edreader

import (
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbn/EDx52display/mfd"

	"github.com/fsnotify/fsnotify"

	"github.com/peterbn/EDx52display/conf"
)

var watcher fsnotify.Watcher

var tick time.Ticker

const (
	pageCargo = iota
	pageLocation
	pageTargetInfo
)

// Mfd is the MFD display structure that will be used by this module. The number of pages should not be changed
var Mfd = mfd.Display{
	Pages: []mfd.Page{
		mfd.Page{
			Lines: []string{},
		},
		mfd.Page{
			Lines: []string{},
		},
		mfd.Page{
			Lines: []string{},
		},
	},
}

// MfdLock locks the current MFD for reads and writes
var MfdLock = sync.RWMutex{}

// PrevMfd is the previous Mfd written to file, to be used for comparisons and avoid superflous updates.
var PrevMfd = Mfd.Copy()

// Start starts the Elite Dangerous journal reader routine
func Start(cfg conf.Conf) {
	// Update immediately, to ensure the mfd.json file exist
	journalfolder := cfg.ExpandJournalFolderPath()
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

	handleCargoFile(filepath.Join(journalfolder, FileCargo))
	swapMfd()
}

// Stop closes the watcher again
func Stop() {
	tick.Stop()
}

func findJournalFile(folder string) string {
	files, _ := filepath.Glob(filepath.Join(folder, "Journal.*.*.log"))
	sort.Strings(files)
	return files[len(files)-1]
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

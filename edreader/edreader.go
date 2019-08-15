package edreader

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/peterbn/EDx52display/mfd"

	"github.com/fsnotify/fsnotify"

	"github.com/peterbn/EDx52display/conf"
)

var watcher fsnotify.Watcher

var tick time.Ticker

const (
	pageCommander = iota
	pageCargo
	pageLocation
)

const (
	commanderHeader = "#     CMDR     #"
	locationHeader  = "#   Location   #"
)

// Mfd is the MFD display structure that will be used by this module. The number of pages should not be changed
var Mfd = mfd.Display{
	Pages: []mfd.Page{
		mfd.Page{
			Lines: []string{commanderHeader},
		},
		mfd.Page{
			Lines: []string{"Cargo: "},
		},
		mfd.Page{
			Lines: []string{locationHeader},
		},
	},
}

// PrevMfd is the previous Mfd written to file, to be used for comparisons and avoid superflous updates.
var PrevMfd = Mfd.Copy()

// Start starts the Elite Dangerous journal reader routine
func Start(cfg conf.Conf) {
	fmt.Println("Config: ", cfg)

	tick := time.NewTicker(time.Duration(cfg.RefreshRateMS) * time.Millisecond)

	go func() {
		for range tick.C {
			// Read in the files at start before we start watching, to initialize
			journalFile := findJournalFile(cfg.JournalsFolder)
			handleJournalFile(journalFile)

			handleCargoFile(filepath.Join(cfg.JournalsFolder, FileCargo))
			swapMfd()
		}
	}()
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
	eq := cmp.Equal(Mfd, PrevMfd)
	if !eq {
		mfd.Write(Mfd)
		PrevMfd = Mfd.Copy()
	}
}

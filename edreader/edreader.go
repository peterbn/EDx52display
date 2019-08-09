package edreader

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/peterbn/EDx52display/mfd"

	"github.com/fsnotify/fsnotify"

	"github.com/peterbn/EDx52display/conf"
)

var watcher fsnotify.Watcher

const (
	pageCargo = iota
)

// Mfd is the MFD display structure that will be used by this module. The number of pages should not be changed
var Mfd = mfd.Display{
	Pages: []mfd.Page{
		mfd.Page{
			Lines: []string{"Cargo: "},
		},
	},
}

// Start starts the Elite Dangerous journal reader routine
func Start(cfg conf.Conf) {
	fmt.Println("Config: ", cfg)

	// Read in the files at start before we start watching, to initialize
	handleCargoFile(path.Join(cfg.JournalsFolder, FileCargo))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				dispatchEvent(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(cfg.JournalsFolder)
	if err != nil {
		log.Fatal(err)
	}
}

// Stop closes the watcher again
func Stop() {
	watcher.Close()
}

func dispatchEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Write == fsnotify.Write {
		log.Println("modified file:", event.Name)
		if isFileEmpty(event.Name) {
			return // don't deal with empty files
		}
		_, file := filepath.Split(event.Name)
		switch file {
		case FileCargo:
			handleCargoFile(event.Name)
		}
	}
}

func isFileEmpty(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
		return true
	}
	return fi.Size() == 0
}

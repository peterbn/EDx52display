package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/peterbn/EDx52display/edreader"
)

func main() {
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		p := edreader.ParseJournalLine(line)
		fmt.Printf("%#v\n", p)
	}

}

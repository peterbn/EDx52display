package main

import (
	"fmt"

	"github.com/peterbn/edx52display/directoutput"
)

func main() {
	directoutput.Initialize("Hello, World")
	defer directoutput.Deinitialize()
	fmt.Println("Initialized program")
	fmt.Scanln()
}

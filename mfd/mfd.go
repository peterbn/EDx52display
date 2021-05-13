package mfd

import (
	"fmt"
)

// Filename is the name of the file written with MFD data
const Filename = "mfd.json"

// Display is the main display structure to write
type Display struct {
	Pages []Page `json:"pages"`
}

// Copy creates a deep copy of this Display
func (d Display) Copy() Display {
	pc := []Page{}
	for _, p := range d.Pages {
		pc = append(pc, p.Copy())
	}
	dc := Display{Pages: pc}

	return dc
}

// Page is a single page of information to show on the MFD
type Page struct {
	Lines []string `json:"lines"`
}

// NewPage returns a new page
func NewPage() Page {
	return Page{Lines: []string{}}
}

// Add appends a new (optionally formatted) string to the LineBuffer
func (p *Page) Add(s string, args ...interface{}) {
	p.Lines = append(p.Lines, fmt.Sprintf(s, args...))
}

// Copy makes a deep copy of this Page
func (p Page) Copy() Page {
	nLines := make([]string, len(p.Lines))
	copy(nLines, p.Lines)
	return Page{Lines: nLines}
}

// Write writes the MFD file
func Write(mfd Display) {
	UpdateDisplay(mfd)
}

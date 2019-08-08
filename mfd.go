package main

// MfdDisplay is the main display structure to write
type MfdDisplay struct {
	Pages []MfdPage `json:"pages"`
}

// MfdPage is a single page of information to show on the MFD
type MfdPage struct {
	Lines []string `json:"lines"`
}

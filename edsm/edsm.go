package edsm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
 Module to communicate with the Elite: Dangerous Star Map site edsm.net
*/

const (
	urlBodies      = "https://www.edsm.net/api-system-v1/bodies?systemId64=%d"
	urlSystemValue = "https://www.edsm.net/api-system-v1/estimated-value?systemId64=%d"
)

// System parses the root object response from the api-system-v1 apis
type System struct {
	ID64      uint64
	Name      string
	BodyCount int

	EstimatedValue       int64
	EstimatedValueMapped int64

	Bodies         []Body
	ValuableBodies []ValuableBody
}

// Body parses information about a single body
type Body struct {
	ID64        uint64
	Name        string
	IsMainStar  bool
	IsScoopable bool
	Type        string
	SubType     string
}

// ValuableBody holds information about the value of bodies
type ValuableBody struct {
	BodyName string
	ValueMax int64
}

// EdsmSystemResult bundles the result of fetching system information with the optional error
type EdsmSystemResult struct {
	S     System
	Error error
}

// MainStar returns the main star in the system
func (s System) MainStar() Body {
	for _, body := range s.Bodies {
		if body.IsMainStar {
			return body
		}
	}
	return Body{}
}

// GetSystemBodies retrieves body information from EDSM.net
func GetSystemBodies(id64 int64) <-chan EdsmSystemResult {
	return getBodyInfo(urlBodies, id64)
}

// GetSystemValue returns information about the system value
func GetSystemValue(id64 int64) <-chan EdsmSystemResult {
	return getBodyInfo(urlSystemValue, id64)
}

var sysinfocache = make(map[string]System)

func getBodyInfo(url string, id64 int64) <-chan EdsmSystemResult {
	retchan := make(chan EdsmSystemResult)
	go func() {
		sysurl := fmt.Sprintf(url, id64)
		cached, ok := sysinfocache[sysurl]
		if ok {
			retchan <- EdsmSystemResult{cached, nil}
			return
		}
		resp, err := http.Get(fmt.Sprintf(url, id64))
		s := System{Bodies: []Body{}}
		if err != nil {
			retchan <- EdsmSystemResult{s, err}
			return
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			retchan <- EdsmSystemResult{s, err}
			return
		}
		json.Unmarshal(data, &s)

		sysinfocache[sysurl] = s
		retchan <- EdsmSystemResult{s, nil}
	}()
	return retchan
}

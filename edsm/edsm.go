package edsm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
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
	ID64   uint64
	BodyID int64

	Name        string
	IsMainStar  bool
	IsScoopable bool
	Type        string
	SubType     string

	Gravity float64

	Volcanism  string
	IsLandable bool

	Materials map[string]float64
}

// ValuableBody holds information about the value of bodies
type ValuableBody struct {
	BodyName string
	ValueMax int64
}

// Material presents a single material and it's presence as a percentage
type Material struct {
	Name       string
	Percentage float64
}

// SystemResult bundles the result of fetching system information with the optional error
type SystemResult struct {
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

// BodyByID retrieves a body from the system by it's BodyID
func (s System) BodyByID(bodyID int64) Body {
	for _, body := range s.Bodies {
		if body.BodyID == bodyID {
			return body
		}
	}
	return Body{}
}

// ShortName returns the shortened name of the body, without the system name prefix
func (b Body) ShortName(s System) string {
	name := b.Name
	if strings.HasPrefix(name, s.Name) {
		return name[len(s.Name)+1:]
	}
	return name
}

// ShortName returns the shortened name of the body, without the system name prefix
func (b ValuableBody) ShortName(s System) string {
	name := b.BodyName
	if strings.HasPrefix(name, s.Name) {
		return name[len(s.Name)+1:]
	}
	return name
}

// MaterialsSorted returns the materials of this body in descending sorted order
func (b Body) MaterialsSorted() []Material {
	ms := []Material{}
	for m, p := range b.Materials {
		ms = append(ms, Material{m, p})
	}

	sort.Slice(ms, func(i, j int) bool {
		if ms[i].Percentage == ms[j].Percentage {
			return ms[i].Name < ms[j].Name
		}
		return ms[i].Percentage > ms[j].Percentage
	})
	return ms
}

// ClearCache will clear the module cache
func ClearCache() {
	cachelock.Lock()
	defer cachelock.Unlock()
	sysinfocache = make(map[string]System)
	log.Debugln("Cached EDSM information cleared")
}

// GetSystemBodies retrieves body information from EDSM.net
func GetSystemBodies(id64 int64) <-chan SystemResult {
	return getBodyInfo(urlBodies, id64)
}

// GetSystemValue returns information about the system value
func GetSystemValue(id64 int64) <-chan SystemResult {
	return getBodyInfo(urlSystemValue, id64)
}

var sysinfocache = make(map[string]System)
var cachelock = sync.RWMutex{}

func getBodyInfo(url string, id64 int64) <-chan SystemResult {
	log.Traceln("getBodyInfo", url, id64)
	retchan := make(chan SystemResult)
	go func() {
		sysurl := fmt.Sprintf(url, id64)

		cachelock.RLock()
		cached, ok := sysinfocache[sysurl]
		cachelock.RUnlock()

		if ok {
			log.Trace("system info found in cache")
			retchan <- SystemResult{cached, nil}
			return
		}
		log.Debugln("Requesting information from EDSM: " + sysurl)
		resp, err := http.Get(fmt.Sprintf(url, id64))
		s := System{Bodies: []Body{}}
		if err != nil {
			retchan <- SystemResult{s, err}
			return
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			retchan <- SystemResult{s, err}
			return
		}
		json.Unmarshal(data, &s)

		cachelock.Lock()
		sysinfocache[sysurl] = s
		cachelock.Unlock()

		retchan <- SystemResult{s, nil}
	}()
	return retchan
}

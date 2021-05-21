package edreader

import (
	"bufio"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/buger/jsonparser"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// LocationType indicates where in a system the player is
type LocationType int

const (
	// LocationSystem means the player is somewhere in the system, not close to a body
	LocationSystem LocationType = iota
	// LocationPlanet means the player is close to a planetary body
	LocationPlanet
	// LocationLanded indicates the player has touched down
	LocationLanded
	// LocationDocked indicates the player has docked at a station (or outpost)
	LocationDocked
)

// Journalstate encapsulates the player state baed on the journal
type Journalstate struct {
	Location
	EDSMTarget
}

// Location indicates the players current location in the game
type Location struct {
	Type LocationType

	SystemAddress int64
	StarSystem    string

	Body     string
	BodyID   int64
	BodyType string

	Latitude  float64
	Longitude float64
}

// EDSMTarget indicates a system targeted by the FSD drive for a jump
type EDSMTarget struct {
	Name          string
	SystemAddress int64
}

const (
	systemaddress = "SystemAddress"
	bodyid        = "BodyID"
	starsystem    = "StarSystem"
	docked        = "Docked"
	body          = "Body"
	bodytype      = "BodyType"
	bodyname      = "BodyName"
	stationname   = "StationName"
	stationtype   = "StationType"
	latitude      = "Latitude"
	longitude     = "Longitude"
	name          = "Name"
)

var state = Journalstate{}

type parser struct {
	line []byte
}

func (p *parser) getString(field string) (string, bool) {
	str, err := jsonparser.GetString(p.line, field)
	if err != nil {
		return "", false
	}
	return str, true
}

func (p *parser) getInt(field string) (int64, bool) {
	num, err := jsonparser.GetInt(p.line, field)
	if err != nil {
		return 0, false
	}
	return num, true
}

func (p *parser) getBool(field string) (bool, bool) {
	b, err := jsonparser.GetBoolean(p.line, field)
	if err != nil {
		return false, false
	}
	return b, true
}

func (p *parser) getFloat(field string) (float64, bool) {
	f, err := jsonparser.GetFloat(p.line, field)
	if err != nil {
		return 0, false
	}
	return f, true
}

var printer = message.NewPrinter(language.English)

// handleJournalFile reads an entire journal file and returns the resulting state
func handleJournalFile(filename string) {
	log.Traceln("Reading journal file " + filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Warnln("Error opening journal file ", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ParseJournalLine(scanner.Bytes())
	}

	RefreshDisplay()
}

// ParseJournalLine parses a single line of the journal and returns the new state after parsing.
func ParseJournalLine(line []byte) Journalstate {
	re := regexp.MustCompile(`"event":"(\w*)"`)
	event := re.FindStringSubmatch(string(line))
	p := parser{line}
	switch event[1] {
	case "Location":
		eLocation(p)
	case "SupercruiseEntry":
		eSupercruiseEntry(p)
	case "SupercruiseExit":
		eSupercruiseExit(p)
	case "FSDJump":
		eFSDJump(p)
	case "Touchdown":
		eTouchDown(p)
	case "Liftoff":
		eLiftoff(p)
	case "FSDTarget":
		eFSDTarget(p)
	case "ApproachBody":
		eApproachBody(p)
	case "ApproachSettlement":
		eApproachSettlement(p)
		break
	}
	return state
}

func eLocation(p parser) {
	// clear current location completely
	state.Type = LocationSystem
	state.Location.SystemAddress, _ = p.getInt(systemaddress)
	state.StarSystem, _ = p.getString(starsystem)

	bodyType, ok := p.getString(bodytype)

	if ok && bodyType == "Planet" {
		state.BodyID, _ = p.getInt(bodyid)
		state.Body, _ = p.getString(body)
		state.BodyType, _ = p.getString(bodytype)
		state.Type = LocationPlanet

		lat, ok := p.getFloat(latitude)
		if ok {
			state.Latitude = lat
			state.Longitude, _ = p.getFloat(longitude)
			state.Type = LocationLanded
		}
	}

	docked, _ := p.getBool(docked)
	if docked {
		state.Type = LocationDocked
	}
}

func eSupercruiseEntry(p parser) {
	state.Type = LocationSystem // don't throw away info
}

func eSupercruiseExit(p parser) {
	eLocation(p)
}

func eFSDJump(p parser) {
	eLocation(p)
}

func eTouchDown(p parser) {
	state.Latitude, _ = p.getFloat(latitude)
	state.Longitude, _ = p.getFloat(longitude)
	state.Type = LocationLanded
}

func eLiftoff(p parser) {
	state.Type = LocationPlanet
}

func eFSDTarget(p parser) {
	state.EDSMTarget.SystemAddress, _ = p.getInt(systemaddress)
	state.EDSMTarget.Name, _ = p.getString(name)

}

func eApproachBody(p parser) {
	state.Location.Body, _ = p.getString(body)
	state.Location.BodyID, _ = p.getInt(bodyid)

	state.Type = LocationPlanet
}

func eApproachSettlement(p parser) {
	state.Location.Body, _ = p.getString(bodyname)
	state.Location.BodyID, _ = p.getInt(bodyid)

	state.Type = LocationPlanet
}

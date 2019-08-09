package journaltypes

// Base is the minimal information for a
type Base struct {
	Event string `json:"event"`
}

// Commander is information about the current commander
type Commander struct {
	Base
	Name string
}

type LoadGame struct {
	Base
	Commander string
	Ship      string
	Credits   int64
}

type Rank struct {
	Base
	Combat, Trade, Explore, Empire, Federation, CQC int
}

type Progress struct {
	Base
	Combat, Trade, Explore, Empire, Federation, CQC int
}

type Reputation struct {
	Base
	Empire, Federation, Alliance float64
}

type Location struct {
	Base
	Docked bool

	StationName       string
	StationType       string
	StationAllegiance string

	StarSystem       string
	SystemSecurity   string `json:"SystemSecurity_Localised"`
	SystemAllegiance string
	SystemFaction    string
	Body             string
	BodyType         string
}

type Docked struct {
	Base

	StationName       string
	StationType       string
	StationAllegiance string
	StarSystem        string
}

type DockingGranted struct {
	Base

	StationName string
	LandingPad  int
}

type FSDJump struct {
	Base

	StarSystem string
}

type Liftoff struct {
	Base
	StarSystem string
}

type SupercruiseEntry struct {
	Base
	StarSystem string
}

type SupercruiseExit struct {
	Base

	Body     string
	BodyType string
}

type Touchdown struct {
	Base

	Latitude  float64
	Longitude float64
}

type Undocked struct {
	Base

	StationName string
}

type ShipyardSwap struct {
	Base

	ShipType string
}

type Continued struct {
	Base
}

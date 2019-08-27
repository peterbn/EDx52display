package edreader

import (
	"bufio"
	"log"
	"os"
	"regexp"

	"github.com/buger/jsonparser"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Journalstate encapsulates the player state baed on the journal
type Journalstate struct {
	Commander string
	Credits   int64
	Ship      string

	Location

	Rank
	Reputation

	EDSMTarget
}

// Rank encapsulates the player's rank
type Rank struct {
	Combat, Trade, Explore, Empire, Federation, CQC string
}

// Reputation indicates the player's reputation with the different factions
type Reputation struct {
	Empire, Federation, Alliance float64
}

// Location indicates the players current location in the game
type Location struct {
	Docked      bool
	Supercruise bool

	StationName       string
	StationType       string
	StationAllegiance string

	SystemAddress    int64
	StarSystem       string
	SystemSecurity   string
	SystemAllegiance string
	SystemFaction    string
	Body             string
	BodyType         string

	Latitude       float64
	Longitude      float64
	HasCoordinates bool
}

// EDSMTarget indicates a system targeted by the FSD drive for a jump
type EDSMTarget struct {
	Name          string
	SystemAddress int64

	BodyTarget bool
	BodyID     int64
}

const (
	name      = "Name"
	commander = "Commander"
	ship      = "Ship"
	credits   = "Credits"
	shiptype  = "ShipType"

	combat     = "Combat"
	trade      = "Trade"
	explore    = "Explore"
	cqc        = "CQC"
	empire     = "Empire"
	federation = "Federation"
	alliance   = "Alliance"

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

	cost             = "Cost"
	totalcost        = "TotalCost"
	basevalue        = "BaseValue"
	bonus            = "Bonus"
	totalsale        = "TotalSale"
	reward           = "Reward"
	totalreward      = "TotalReward"
	transfercost     = "TransferCost"
	donation         = "Donation"
	buyprice         = "BuyPrice"
	sellprice        = "SellPrice"
	amount           = "Amount"
	shipprice        = "ShipPrice"
	transferprice    = "TransferPrice"
	totalearnings    = "TotalEarnings"
	brokerpercentage = "BrokerPercentage"
)

var state = Journalstate{}

func (s *Journalstate) addCredits(c int64) {
	s.Credits = s.Credits + c
}

func (s *Journalstate) subCredits(c int64) {
	s.Credits = s.Credits - c
}

type parser struct {
	line []byte
}

func (p *parser) getString(field string) string {
	str, err := jsonparser.GetString(p.line, field)
	if err != nil {
		return ""
	}
	return str
}

func (p *parser) getInt(field string) int64 {
	num, err := jsonparser.GetInt(p.line, field)
	if err != nil {
		return 0
	}
	return num
}

func (p *parser) getBool(field string) bool {
	b, err := jsonparser.GetBoolean(p.line, field)
	if err != nil {
		return false
	}
	return b
}

func (p *parser) getFloat(field string) float64 {
	f, err := jsonparser.GetFloat(p.line, field)
	if err != nil {
		return 0
	}
	return f
}

var printer = message.NewPrinter(language.English)

// handleJournalFile reads an entire journal file and returns the resulting state
func handleJournalFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening journal file ", filename, err)
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
	case "LoadGame":
		eLoadGame(p)
	case "Commander":
		eCommander(p)
	case "Rank":
		eRank(p)
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
	case "Docked":
		eDocked(p)
	case "Undocked":
		eUndocked(p)
	case "BuyExplorationData":
		costEvent(p, cost)
	case "SellExplorationData":
		gainEvent(p, totalearnings)
	case "BuyTradeData":
		costEvent(p, cost)
	case "MarketBuy":
		costEvent(p, totalcost)
	case "MarketSell":
		gainEvent(p, totalsale)
	case "BuyAmmo":
		costEvent(p, cost)
	case "BuyDrones":
		costEvent(p, totalcost)
	case "CommunityGoalReward":
		gainEvent(p, reward)
	case "CrewHire":
		costEvent(p, cost)
	case "FetchRemoteModule":
		costEvent(p, transfercost)
	case "MissionCompleted":
		eMissionCompleted(p)
	case "ModuleBuy":
		eModuleBuy(p)
	case "ModuleSell":
		gainEvent(p, sellprice)
	case "ModuleSellRemote":
		gainEvent(p, sellprice)
	case "PayFines":
		costEvent(p, amount)
	case "PayLegacyFines":
		costEvent(p, amount)
	case "RedeemVoucher":
		eRedeemVoucher(p)
	case "RefuelAll":
		costEvent(p, cost)
	case "RefuelPartial":
		costEvent(p, cost)
	case "Repair":
		costEvent(p, cost)
	case "RepairAll":
		costEvent(p, cost)
	case "RestockVehicle":
		costEvent(p, cost)
	case "SellDrones":
		gainEvent(p, totalsale)
	case "ShipyardBuy":
		eShipyardBuy(p)
	case "ShipyardSell":
		gainEvent(p, shipprice)
	case "ShipyardTransfer":
		costEvent(p, transferprice)
	case "ShipyardSwap":
		eShipyardSwap(p)
	case "PowerplayFastTrack":
		costEvent(p, cost)
	case "PowerplaySalary":
		gainEvent(p, amount)
	case "Bounty":
		gainEvent(p, totalreward)
	case "DatalinkVoucher":
		gainEvent(p, reward)
	case "FactionKillBond":
		gainEvent(p, reward)
	case "MultiSellExplorationData":
		gainEvent(p, totalearnings)
	case "PayBounties":
		costEvent(p, amount)
	case "Promotion":
		ePromotion(p)
	case "Reputation":
		eReputation(p)
	case "Liftoff":
		eLiftoff(p)
	case "FSDTarget":
		eFSDTarget(p)
	case "ApproachBody":
		eApproachBody(p)
	case "ApproachSettlement":
		eApproachSettlement(p)
	case "AfmuRepairs",
		"AsteroidCracked",
		"Cargo",
		"CargoDepot",
		"CodexEntry",
		"CollectCargo",
		"CommitCrime",
		"CommunityGoal",
		"CommunityGoalJoin",
		"CrimeVictim",
		"DatalinkScan",
		"DataScanned",
		"Died",
		"DiscoveryScan",
		"DockingDenied",
		"DockingGranted",
		"DockingRequested",
		"DockingTimeout",
		"DockSRV",
		"EjectCargo",
		"EngineerApply",
		"EngineerContribution",
		"EngineerCraft",
		"EngineerProgress",
		"EscapeInterdiction",
		"Fileheader",
		"FSSAllBodiesFound",
		"FSSDiscoveryScan",
		"FSSSignalDiscovered",
		"FuelScoop",
		"HeatDamage",
		"HeatWarning",
		"HullDamage",
		"Interdicted",
		"JetConeBoost",
		"LaunchDrone",
		"LaunchSRV",
		"LeaveBody",
		"Loadout",
		"Market",
		"MaterialCollected",
		"MaterialDiscovered",
		"Materials",
		"MaterialTrade",
		"MiningRefined",
		"MissionAbandoned",
		"MissionAccepted",
		"MissionFailed",
		"MissionRedirected",
		"Missions",
		"ModuleInfo",
		"ModuleRetrieve",
		"ModuleStore",
		"ModuleSwap",
		"Music",
		"NavBeaconScan",
		"NewCommander",
		"Outfitting",
		"Passengers",
		"Progress",
		"ProspectedAsteroid",
		"ReceiveText",
		"ReservoirReplenished",
		"Resurrect",
		"SAAScanComplete",
		"Scan",
		"Scanned",
		"Screenshot",
		"SearchAndRescue",
		"SendText",
		"SetUserShipName",
		"ShieldState",
		"ShipTargeted",
		"Shipyard",
		"ShipyardNew",
		"Shutdown",
		"StartJump",
		"Statistics",
		"StoredModules",
		"StoredShips",
		"Synthesis",
		"SystemsShutdown",
		"TechnologyBroker",
		"UnderAttack",
		"USSDrop":
		break
	}
	return state
}

func eCommander(p parser) {
	state.Commander = p.getString(name)
}

func eLoadGame(p parser) {
	state.Commander = p.getString(commander)
	state.Ship = p.getString(ship)
	state.Credits = p.getInt(credits)
}

func eRank(p parser) {
	state.Combat = combatRank[p.getInt(combat)]
	state.Trade = tradeRank[p.getInt(trade)]
	state.Explore = explorerRank[p.getInt(explore)]
	state.CQC = cqcRank[p.getInt(cqc)]

	state.Rank.Empire = empireRank[p.getInt(empire)]
	state.Rank.Federation = federationRank[p.getInt(federation)]
}

func eLocation(p parser) {
	// clear current location completely
	state.Location = Location{}
	state.Location.SystemAddress = p.getInt(systemaddress)
	state.StarSystem = p.getString(starsystem)
	state.Docked = p.getBool(docked)

	state.Body = p.getString(body)
	state.BodyType = p.getString(bodytype)

	state.StationName = p.getString(stationname)
	state.StationType = p.getString(stationtype)
}

func eSupercruiseEntry(p parser) {
	eLocation(p)
	state.Supercruise = true

	if state.EDSMTarget.BodyTarget {
		state.EDSMTarget.BodyTarget = false
		state.EDSMTarget.SystemAddress = 0
	}
}

func eSupercruiseExit(p parser) {
	eLocation(p)

	if p.getString(bodytype) == "Planet" && !state.EDSMTarget.BodyTarget {
		// If we don't already have bodyInfo from an approach event
		eApproachBody(p)
	}
}

func eFSDJump(p parser) {
	eLocation(p)
	state.EDSMTarget = EDSMTarget{}
	state.Supercruise = true
}

func eTouchDown(p parser) {
	state.Latitude = p.getFloat(latitude)
	state.Longitude = p.getFloat(longitude)
	state.HasCoordinates = true
}

func eLiftoff(p parser) {
	state.HasCoordinates = false
}

func eDocked(p parser) {
	eLocation(p)
	state.Docked = true
}

func eUndocked(p parser) {
	eLocation(p)
}

func costEvent(p parser, key string) {
	c := p.getInt(key)
	state.subCredits(c)
}

func gainEvent(p parser, key string) {
	g := p.getInt(key)
	state.addCredits(g)
}

func eMissionCompleted(p parser) {
	r := p.getInt(reward)
	d := p.getInt(donation)
	state.addCredits(r - d)
}

func eModuleBuy(p parser) {
	buy := p.getInt(buyprice)
	sell := p.getInt(sellprice)

	// Any optional sale price is positive, buy price is negative
	state.addCredits(sell - buy)
}

func eShipyardBuy(p parser) {
	price := p.getInt(shipprice)
	sale := p.getInt(sellprice)

	state.addCredits(sale - price)
	state.Ship = p.getString(shiptype)
}

func eShipyardSwap(p parser) {
	state.Ship = p.getString(shiptype)
}

func eRedeemVoucher(p parser) {
	total := p.getInt(amount)
	fee := p.getFloat(brokerpercentage)

	if fee > 0 {
		total = (total * (100 - int64(fee))) / 100
	}
	state.addCredits(total)
}

func ePromotion(p parser) {
	state.Combat = combatRank[p.getInt(combat)]
	state.Trade = tradeRank[p.getInt(trade)]
	state.Explore = explorerRank[p.getInt(explore)]
	state.CQC = cqcRank[p.getInt(cqc)]
}

func eReputation(p parser) {
	state.Reputation.Empire = p.getFloat(empire)
	state.Reputation.Federation = p.getFloat(federation)
	state.Reputation.Alliance = p.getFloat(alliance)
}

func eFSDTarget(p parser) {
	state.EDSMTarget.SystemAddress = p.getInt(systemaddress)
	state.EDSMTarget.Name = p.getString(name)
	state.EDSMTarget.BodyTarget = false

}

func eApproachBody(p parser) {
	state.EDSMTarget.SystemAddress = p.getInt(systemaddress)
	state.EDSMTarget.Name = p.getString(body)

	state.EDSMTarget.BodyTarget = true
	state.EDSMTarget.BodyID = p.getInt(bodyid)
}

func eApproachSettlement(p parser) {
	state.EDSMTarget.SystemAddress = p.getInt(systemaddress)
	state.EDSMTarget.Name = p.getString(bodyname)

	state.EDSMTarget.BodyTarget = true
	state.EDSMTarget.BodyID = p.getInt(bodyid)
}

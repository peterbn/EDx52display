package edreader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/buger/jsonparser"
	"github.com/peterbn/EDx52display/mfd"
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

	starsystem  = "StarSystem"
	docked      = "Docked"
	body        = "Body"
	bodytype    = "BodyType"
	stationname = "StationName"
	stationtype = "StationType"
	latitude    = "Latitude"
	longitude   = "Longitude"

	cost          = "Cost"
	totalcost     = "TotalCost"
	basevalue     = "BaseValue"
	bonus         = "Bonus"
	totalsale     = "TotalSale"
	reward        = "Reward"
	totalreward   = "TotalReward"
	transfercost  = "TransferCost"
	donation      = "Donation"
	buyprice      = "BuyPrice"
	sellprice     = "SellPrice"
	amount        = "Amount"
	shipprice     = "ShipPrice"
	transferprice = "TransferPrice"
	totalearnings = "TotalEarnings"
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

	cmdr := []string{commanderHeader}

	printer := message.NewPrinter(language.English)

	cmdr = append(cmdr, state.Commander)
	cmdr = append(cmdr, printer.Sprintf("CR:%d", state.Credits))
	cmdr = append(cmdr, state.Ship)

	loc := []string{locationHeader}
	loc = append(loc, state.StarSystem)
	if state.StationName != "" {
		loc = append(loc, state.StationName)
	}
	if len(state.StationType) > 0 {
		loc = append(loc, state.StationType)
	}
	if state.Body != "" {
		loc = append(loc, state.Body)
	}
	if state.BodyType != "" {
		loc = append(loc, state.BodyType)
	}
	if state.HasCoordinates {
		loc = append(loc, fmt.Sprintf("Lat: %.3f", state.Latitude))
		loc = append(loc, fmt.Sprintf("Lon: %.3f", state.Longitude))
	}

	if state.Supercruise {
		loc = append(loc, "In Supercruise")
	}
	if state.Docked {
		loc = append(loc, "Docked")
	}
	Mfd.Pages[pageCommander] = mfd.Page{Lines: cmdr}
	Mfd.Pages[pageLocation] = mfd.Page{Lines: loc}
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
		gainEvent(p, amount)
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
	case "AfmuRepairs",
		"ApproachBody",
		"ApproachSettlement",
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
 		"FSDTarget", // Totally gonna parse this and request information from edsm.net o_O
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
	state.StarSystem = p.getString(starsystem)
	state.Docked = p.getBool(docked)

	state.Body = p.getString(body)
	state.BodyType = p.getString(bodytype)

	state.StationName = p.getString(stationname)
	state.StationType = p.getString(stationtype)
}

func eSupercruiseEntry(p parser) {
	state.Location = Location{}

	state.Supercruise = true
	state.StarSystem = p.getString(starsystem)
}

func eSupercruiseExit(p parser) {
	eLocation(p)
}

func eFSDJump(p parser) {
	eLocation(p)
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
	state.Docked = false
	state.StationName = p.getString(stationname)
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

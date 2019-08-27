package edreader

import (
	"fmt"
	"log"

	"github.com/peterbn/EDx52display/edsm"
	"github.com/peterbn/EDx52display/mfd"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// RefreshDisplay updates the display with the current state
func RefreshDisplay() {
	MfdLock.Lock()
	defer MfdLock.Unlock()
	Mfd.Pages[pageCommander] = mfd.Page{Lines: renderCmdrPage()}
	Mfd.Pages[pageLocation] = mfd.Page{Lines: renderLocationPage()}
	Mfd.Pages[pageSysInfo] = mfd.Page{Lines: renderEDSMData()}
}

func renderCmdrPage() []string {
	cmdr := []string{commanderHeader}

	printer := message.NewPrinter(language.English)

	cmdr = append(cmdr, state.Commander)
	cmdr = append(cmdr, "Credits:")
	cmdr = append(cmdr, printer.Sprintf("%16d", state.Credits))

	cmdr = append(cmdr, "Combat:")
	cmdr = append(cmdr, fmt.Sprintf("%16s", state.Combat))
	cmdr = append(cmdr, "Trade:")
	cmdr = append(cmdr, fmt.Sprintf("%16s", state.Trade))
	cmdr = append(cmdr, "Exploration:")
	cmdr = append(cmdr, fmt.Sprintf("%16s", state.Explore))
	cmdr = append(cmdr, "CQC:")
	cmdr = append(cmdr, fmt.Sprintf("%16s", state.CQC))

	cmdr = append(cmdr, fmt.Sprintf("Empire:%9.1f", state.Reputation.Empire))
	cmdr = append(cmdr, fmt.Sprintf("%16s", state.Rank.Empire))

	cmdr = append(cmdr, fmt.Sprintf("Federation:%5.1f", state.Reputation.Federation))
	cmdr = append(cmdr, fmt.Sprintf("%16s", state.Rank.Federation))

	cmdr = append(cmdr, fmt.Sprintf("Alliance:%7.1f", state.Reputation.Alliance))
	return cmdr
}

func renderLocationPage() []string {
	loc := []string{locationHeader}
	if state.Supercruise {
		loc = append(loc, "In Supercruise")
	}
	if state.Docked {
		loc = append(loc, "Docked")
	}
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
	return loc
}

func renderEDSMData() []string {
	if state.EDSMTarget.BodyTarget {
		return renderEDSMBody("#     Body     #", state.EDSMTarget.Name, state.EDSMTarget.SystemAddress, state.EDSMTarget.BodyID)
	} else if state.EDSMTarget.SystemAddress != 0 {
		return renderEDSMSystem("#  System <T>  #", state.EDSMTarget.Name, state.EDSMTarget.SystemAddress)
	} else {
		return renderEDSMSystem("#  System <C>  #", state.Location.StarSystem, state.Location.SystemAddress)
	}
}

func renderEDSMSystem(header, systemname string, systemaddress int64) []string {
	sysinfopromise := edsm.GetSystemBodies(systemaddress)
	valueinfopromise := edsm.GetSystemValue(systemaddress)

	sys := []string{header, systemname}

	sysinfo := <-sysinfopromise

	if sysinfo.Error != nil {
		log.Println("Unable to fetch system information: ", sysinfo.Error)
		sys = append(sys, "Sysinfo lookup error")
	} else if sysinfo.S.ID64 == 0 {
		sys = append(sys, "No EDSM data")
	} else {
		mainBody := sysinfo.S.MainStar()
		if mainBody.IsScoopable {
			sys = append(sys, "Scoopable")
		} else {
			sys = append(sys, "Not scoopable")
		}

		sys = append(sys, mainBody.SubType)

		sys = append(sys, fmt.Sprintf("Bodies: %d", sysinfo.S.BodyCount))

		valinfo := <-valueinfopromise
		if valinfo.Error == nil {
			sys = append(sys, "Scan value:")
			sys = append(sys, printer.Sprintf("%16d", valinfo.S.EstimatedValue))
			sys = append(sys, "Mapped value:")
			sys = append(sys, printer.Sprintf("%16d", valinfo.S.EstimatedValueMapped))

			if len(valinfo.S.ValuableBodies) > 0 {
				sys = append(sys, "Valuable Bodies:")
			}
			for _, valbody := range valinfo.S.ValuableBodies {
				sys = append(sys, valbody.BodyName)
				sys = append(sys, printer.Sprintf("%16d", valbody.ValueMax))
			}

		}
	}
	return sys
}

func renderEDSMBody(header, bodyName string, systemaddress, bodyid int64) []string {
	sysinfopromise := edsm.GetSystemBodies(systemaddress)
	page := []string{header, bodyName}
	sysinfo := <-sysinfopromise
	if sysinfo.Error != nil {
		log.Println("Unable to fetch system information: ", sysinfo.Error)
		page = append(page, "Sysinfo lookup error")
		return page
	}
	if sysinfo.S.ID64 == 0 {
		page = append(page, "No EDSM data")
		return page
	}

	body := sysinfo.S.BodyByID(bodyid)
	if body.BodyID == 0 {
		page = append(page, "No EDSM data")
		return page
	}

	page = append(page, fmt.Sprintf("Gravity %7.2fG", body.Gravity))

	page = append(page, "Materials:")
	for _, m := range body.MaterialsSorted() {
		page = append(page, fmt.Sprintf("%5.2f%% %s", m.Percentage, m.Name))
	}

	return page
}

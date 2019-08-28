package edreader

import (
	"fmt"
	"log"

	"github.com/peterbn/EDx52display/edsm"
	"github.com/peterbn/EDx52display/mfd"
)

// RefreshDisplay updates the display with the current state
func RefreshDisplay() {
	MfdLock.Lock()
	defer MfdLock.Unlock()
	Mfd.Pages[pageLocation] = mfd.Page{Lines: renderLocationPage()}
	Mfd.Pages[pageTargetInfo] = mfd.Page{Lines: renderFSDTarget()}
}

func renderLocationPage() []string {
	if state.Type == LocationPlanet || state.Type == LocationLanded {
		return renderEDSMBody("#     Body     #", state.Location.Body, state.Location.SystemAddress, state.BodyID)
	}

	return renderEDSMSystem("#    System    #", state.Location.StarSystem, state.Location.SystemAddress)
}

func renderFSDTarget() []string {
	if state.EDSMTarget.SystemAddress == 0 {
		return []string{"No FSD Target"}
	}
	return renderEDSMSystem("#  FSD Target  #", state.EDSMTarget.Name, state.EDSMTarget.SystemAddress)
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

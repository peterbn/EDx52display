# EDx52display

Reading Elite: Dangerous journal information and displaying on a Logitech X52 PRO MFD via [my JSON-based driver](https://github.com/peterbn/X52-pro-MFD-JSON-Driver).

    NOTE: It is recommended to run a tool that uploads data to EDSM, such as [ED Market Connector](https://github.com/Marginal/EDMarketConnector). <br>
    Doing this will ensure that any new discoveries can be shown on the display.

## Installation

Unzip the release folder and run the `EDx52display.exe` application.

## Output

Running this application will show 3 pages of information on your MFD. Most of this information is sourced from EDSM.net.

Of particular note is:

- Live view of cargo hold
- Value of scanning and mapping the system
- Surface gravity of the planet you are about to land on

### Page 1: Cargo hold

This page will simply show the total used capacity and the contents of your cargo hold. This can be useful when mining, to check progress without having to go into the inventory panel.

### Page 2: Current location

This page will show information about your current location.
Currently, this means either the system you are in, or the planet you have approached.
See below for specifics of what is shown for each type

### Page 3: FSD Target

This page will show system information about the system targeted for a FSD Jump

### System Page

A page with system information will have the following information, sourced from EDSM:

- System Name
- Whether the main star is scoopable
- Number of bodies (as reported by EDSM)
- Total value for scanning the system
- Total value for mapping the entire system
- Any valuable bodies
- System Prospecting information
- - Available elements, with number of planets landable where they occur
- - The planet in the system with the highest occurence of said element

### Planet Page

A page with planet information will have the following data, sourced from EDSM:

- Planet name
- Planet Gravity (!)
- Available materials for the planet, if any

## Troubleshooting

This application reads the journal files of your elite dangerous installation.
These are normally located in `%USERPROFILE%\\Saved Games\\Frontier Developments\\Elite Dangerous` on Windows. However, if your installation
uses a different location, you should update the conf.yaml file in the installation folder.

## Credits

This project owes a great deal to [Anthony Zaprzalka](https://github.com/AZaps) in terms of idea and execution
and to [Jonathan Harris](https://github.com/Marginal) and the [EDMarketConnector](https://github.com/Marginal/EDMarketConnector) project
for the CSV files of names for all the commodities.

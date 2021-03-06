package main

var usage = `Usage:
	dashboard-client -c config_file
	dashboard-client -version
	dashbaord-client {-update|-checkupdate}

DESCRIPTION:
	dashboard-client periodically sends information to an API end point.

OPTIONS:
	-c config_file:
		The configuration file (see CONFIGURAITON FILE)

	-version:
		Display version information

	-update:
		Update the binary

	-checkupdate
		Check if updates are available for the binary

CONFIGURATION FILE:
	The configuration file is formatted in JSON and has the following layout:

	{
		"api" : "http://dashboard.example.org/api",
		"identifier": "webserver",
		"name": "Web Server",
		"widgets" : [
			"uptime" : {},
			"load" : {},
			"time" : {},
			"network" : {
				"interface" : "eth0"
			}
		]
	}	

	**api**: The dashboard server API end point.
	**identifier**: Unique identifier for the client.
	**name**: Human-readable name of the client.
	**widgets**: Widgets and configuration for which to send information.
`

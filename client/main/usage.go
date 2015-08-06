package main

var usage = `Usage: dashboard-client -c config_file

DESCRIPTION:
	dashboard-client periodically sends information to an API end point.

OPTIONS:
	-c config_file:
		The configuration file (see CONFIGURAITON FILE)

CONFIGURATION FILE:
	The configuration file is formatted in JSON and has the following layout:

	{
		"api" : "http://dashboard.example.org/api",
		"widgets" : [
			"uptime",
			"load",
			"time"
		]
	}	

	**api**: The dashboard server API end point.
	**widgets**: Widgets for which to send information.
`

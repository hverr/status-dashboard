package main

var usage = `Usage: dashboard-server -c config_file

DESCRIPTION:
	dashboard-server servers the server application

OPTIONS:
	-c config_file:
		The configuration file (see CONFIGURAITON FILE)

CONFIGURATION FILE:
	The configuration file is formatted in JSON and has the following layout:

	{
		"clients" : [
			"webserver",
			"mysql"
		],
		"defaultWidgets" : {
			"webserver" : [
				"uptime",
				"load"
			],
			"mysql" : [
				"load",
				"queries"
			]
		}
	}

	**clients**: List of allowed clients.
	**defaultWidgets**: Object mapping a client to its default widgets.
`

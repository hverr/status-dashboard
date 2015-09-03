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
		"updateInterval" : 3,
		"clients" : {
			"webserver" : {
				"secret" : "supersecretkey"
			},
			"mysql" : {}
		}
	}

	**updateInterval**: Minimum time between update requests. (optional)
	**clients**: List of allowed clients and their configuration.
`

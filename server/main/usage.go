package main

var usage = `Usage:
	dashboard-server -c config_file [-listen listen_addr] [-debug]
	dashboard-server -version
	dashbaord-server {-update|-checkupdate}

DESCRIPTION:
	dashboard-server servers the server application

OPTIONS:
	-c config_file:
		The configuration file (see CONFIGURAITON FILE)

	-listen listen_addr:
		The address to bind to

	-debug:
		Launch the HTTP router in debug mode

	-version:
		Show version information.

	-update:
		Update the binary

	-checkupdate
		Check if updates are available for the binary

CONFIGURATION FILE:
	The configuration file is formatted in JSON and has the following layout:

	{
		"updateInterval" : 3,
		"clients" : {
			"webserver" : {
				"secret" : "supersecretkey"
			},
			"mysql" : {}
		},
		"users" : {
		  "username" : "password"
		}
	}

	**updateInterval**: Minimum time between update requests. (optional)
	**clients**: List of allowed clients and their configuration.
	**users**: Dictionary of login:password pairs for users that are allowed to access the dashboard. (optional)
`

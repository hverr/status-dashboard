package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/api"
	"github.com/hverr/status-dashboard/server/scheduler"
	"github.com/hverr/status-dashboard/server/static"
	"github.com/hverr/status-dashboard/version"
)

func main() {
	var configFile string
	var listenAddr string
	var debug bool
	var showVersion bool

	flag.StringVar(&configFile, "c", "", "Configuration file.")
	flag.StringVar(&listenAddr, "listen", ":8050", "Listen address.")
	flag.BoolVar(&debug, "debug", false, "Debug gin router.")
	flag.BoolVar(&showVersion, "version", false, "Show version information.")
	flag.Parse()

	printHelp := func() {
		fmt.Fprint(os.Stderr, usage)
	}

	if showVersion {
		if !version.PrintVersionInformation(os.Stdout) {
			os.Exit(1)
		}
		return
	}

	if configFile == "" {
		fmt.Fprintln(os.Stderr, "fatal: flag: missing -c option")
		fmt.Fprintln(os.Stderr, "")
		printHelp()
		os.Exit(1)
	}

	// Setup configuration
	config := server.Configuration{}
	if err := config.ParseConfiguration(configFile); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not parse configuration file ",
			configFile+":", err)
		os.Exit(1)
	}

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Setup user authenticatuer
	userAuth := server.UserAuthenticator{
		Configuration: config,
	}

	// Create a server
	serv := server.NewServer(config)

	// Install static api
	staticApi := static.Static{
		UserAuthenticator: userAuth,
	}
	if err := staticApi.Install(router); err != nil {
		log.Println("We could not serve the static files. If not already")
		log.Println("done, set the HTML_ROOT environment variable to the")
		log.Println("of the static files.")
		log.Fatal("fatal: could not serve static files:", err)
	}

	// Install rest api
	restApi := api.API{
		Configuration:     config,
		Server:            serv,
		UserAuthenticator: userAuth,
		Scheduler:         scheduler.New(config),
	}
	if err := restApi.Install(router); err != nil {
		log.Fatal("fatal: could not serve API:", err)
	}

	router.Run(listenAddr)
}

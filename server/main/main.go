package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/settings"
	"github.com/hverr/status-dashboard/server/static"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "c", "", "Configuration file.")
	flag.Parse()

	printHelp := func() {
		fmt.Fprint(os.Stderr, usage)
	}

	if configFile == "" {
		fmt.Fprintln(os.Stderr, "fatal: flag: missing -c option")
		fmt.Fprintln(os.Stderr, "")
		printHelp()
		os.Exit(1)
	}

	if err := server.ParseConfiguration(configFile); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not parse configuration file ",
			configFile+":", err)
		os.Exit(1)
	}

	router := gin.Default()

	if err := static.Install(router); err != nil {
		log.Println("fatal: could not serve static files:", err)
		os.Exit(1)
	}

	router.Run(settings.ListenAddress)
}

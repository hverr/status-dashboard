package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/api"
	"github.com/hverr/status-dashboard/server/static"
)

func main() {
	var configFile string
	var listenAddr string

	flag.StringVar(&configFile, "c", "", "Configuration file.")
	flag.StringVar(&listenAddr, "listen", ":8050", "Listen address.")
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
		log.Println("We could not serve the static files. If not already")
		log.Println("done, set the HTML_ROOT environment variable to the")
		log.Println("of the static files.")
		log.Fatal("fatal: could not serve static files:", err)
	}

	if err := api.Install(router); err != nil {
		log.Fatal("fatal: could not serve API:", err)
	}

	router.Run(listenAddr)
}

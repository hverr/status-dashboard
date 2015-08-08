package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hverr/status-dashboard/client"
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

	if err := client.ParseConfiguration(configFile); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not parse configuration file",
			configFile+":", err)
		os.Exit(1)
	}

	fmt.Println("Using config:", client.Configuration)
}

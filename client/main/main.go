package main

import (
	"flag"
	"fmt"
	"os"
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
		return
	}

	fmt.Println("Should parse ", configFile)
}

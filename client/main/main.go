package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hverr/status-dashboard/client"
	"github.com/hverr/status-dashboard/widgets"
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

	if err := client.Register(); err != nil {
		log.Fatal(err)
	}

	for {
		if err := update(); err != nil {
			log.Println("Could not send updates:", err)
		} else {
			log.Println("Sent widget information.")
		}
	}
}

func update() error {
	requested, err := client.GetRequestedWidgets()
	if err != nil {
		return err
	}

	results := make([]widgets.Widget, 0, len(requested.Widgets))

	for _, w := range requested.Widgets {
		initiator := widgets.AllWidgets[w]
		if initiator == nil {
			return errors.New("Unknown requested widget type: " + w)
		}

		widget := initiator()
		if err := widget.Update(); err != nil {
			return err
		}

		results = append(results, widget)
	}

	for _, w := range results {
		if err := client.PostWidgetUpdate(w); err != nil {
			return err
		}
	}

	return nil
}

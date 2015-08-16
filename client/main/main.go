package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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

	for {
		registerUpdateLoop()

		t := 5 * time.Second
		log.Println("Reregistering in", t)
		<-time.After(t)
	}
}

func registerUpdateLoop() {
	if err := client.Register(); err != nil {
		log.Println(err)
		return
	}

	for {
		if err := update(); err != nil {
			log.Println("Could not send updates:", err)
			return
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

	results := make([]widgets.BulkElement, 0, len(requested.Widgets))

	for _, w := range requested.Widgets {
		var widget widgets.Widget

		initiator := widgets.AllWidgets[w]
		if initiator == nil {
			log.Print("Unknown requested widget type: " + w)
			widget = nil

		} else {
			widget = initiator()
			if err := widget.Update(); err != nil {
				log.Printf("Can't update %v: %v", w, err)
				widget = nil
			}
		}

		e := widgets.BulkElement{
			Type:   w,
			Widget: widget,
		}

		results = append(results, e)
	}

	if err := client.PostWidgetBulkUpdate(results); err != nil {
		return err
	}

	return nil
}

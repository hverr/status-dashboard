package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hverr/status-dashboard/client"
	"github.com/hverr/status-dashboard/widgets"
)

func main() {
	var configFile string
	var ca string

	flag.StringVar(&configFile, "c", "", "Configuration file.")
	flag.StringVar(&ca, "ca", "", "Root CA certificate.")
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

	if ca != "" {
		var config tls.Config
		pem, err := ioutil.ReadFile(ca)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: read %v: %v\n", ca, err)
			os.Exit(1)
		}

		config.RootCAs = x509.NewCertPool()
		if !config.RootCAs.AppendCertsFromPEM(pem) {
			log.Println("warning: x509: could not use PEM in", ca)
		}

		client.Session.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &config,
			},
		}
	}

	if err := client.ParseConfiguration(configFile); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not parse configuration file",
			configFile+":", err)
		os.Exit(1)
	}

	all, err := initializeWidgets()
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not initialize widgets:", err)
		os.Exit(1)
	}

	for {
		registerUpdateLoop(all)

		t := 5 * time.Second
		log.Println("Reregistering in", t)
		<-time.After(t)
	}
}

func initializeWidgets() ([]widgets.Widget, error) {
	all := make([]widgets.Widget, 0)
	for identifier, configuration := range client.Configuration.Widgets {
		initiator := widgets.AllWidgets[identifier]
		if initiator == nil {
			return nil, errors.New("Unsupported widget " + identifier)
		}

		w := initiator()
		if err := w.Configure(configuration); err != nil {
			return nil, fmt.Errorf("Could not configure widget %s: %v", identifier, err)
		}
		if err := w.Start(); err != nil {
			return nil, fmt.Errorf("Could not start widget %s: %v", identifier, err)
		}

		all = append(all, w)
	}

	return all, nil
}

func registerUpdateLoop(allWidgets []widgets.Widget) {
	if err := client.Register(allWidgets); err != nil {
		log.Println(err)
		return
	}

	for {
		if err := update(allWidgets); err != nil {
			log.Println("Could not send updates:", err)
			return
		} else {
			log.Println("Sent widget information.")
		}
	}
}

func update(allWidgets []widgets.Widget) error {
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
			if widget.HasData() == false {
				widget = nil
			} else if err := widget.Update(); err != nil {
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

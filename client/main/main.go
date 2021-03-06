package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hverr/status-dashboard/client"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/updates"
	"github.com/hverr/status-dashboard/version"
	"github.com/hverr/status-dashboard/widgets"
)

const BinaryBaseName = "dashboard-client"

type environment struct {
	Server        client.Server
	Configuration client.Configuration
}

func main() {
	var configFile string
	var ca string
	var showVersion bool
	var shouldUpdate bool
	var checkUpdate bool

	flag.StringVar(&configFile, "c", "", "Configuration file.")
	flag.StringVar(&ca, "ca", "", "Root CA certificate.")
	flag.BoolVar(&showVersion, "version", false, "Show version information.")
	flag.BoolVar(&shouldUpdate, "update", false, "Update the binary.")
	flag.BoolVar(&checkUpdate, "checkupdate", false, "Check for updates.")
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

	if checkUpdate {
		updates.CheckForUpdates(BinaryBaseName)
		return
	}

	if shouldUpdate {
		updates.UpdateApp(BinaryBaseName)
		return
	}

	if configFile == "" {
		fmt.Fprintln(os.Stderr, "fatal: flag: missing -c option")
		fmt.Fprintln(os.Stderr, "")
		printHelp()
		os.Exit(1)
	}

	env := &environment{}

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

		env.Server.Session.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &config,
			},
		}
	}

	// Parse configuration file.
	fh, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	if err := env.Configuration.ParseConfiguration(fh); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not parse configuration file",
			configFile+":", err)
		os.Exit(1)
	}

	// Configure server behaviour
	env.Server.Configuration = env.Configuration

	for {
		// Register the client.
		if initialized, err, recoverable := register(env); err != nil {
			if !recoverable {
				log.Fatal("Could not register:", err)
			} else {
				log.Println("Could not register:", err)
			}
		} else {
			// Registration was successful, constantly update now.
			log.Println("Successfully registered:", initialized)
			started := make(map[string]widgets.Widget)
			for {
				if err := update(env, initialized, started); err != nil {
					log.Println("Could not update widgets:", err)
					break
				}
			}
		}

		// Don't overload
		t := 5 * time.Second
		log.Println("Reregistering in", t)
		<-time.After(t)
	}
}

func register(env *environment) (initialized map[string]widgets.Widget, err error, recoverable bool) {
	// Determine available widgets and initialize them.
	initialized = make(map[string]widgets.Widget)
	availableWidgets := make([]server.WidgetRegistration, 0)
	for widgetType, config := range env.Configuration.Widgets {
		initiator := widgets.AllWidgets[widgetType]
		if initiator == nil {
			return nil, fmt.Errorf("Unsupported widget " + widgetType), false
		}

		w := initiator()
		if err := w.Configure(config); err != nil {
			return nil, fmt.Errorf("Could not configure %s: %v", widgetType, err), false
		}

		r := server.WidgetRegistration{
			Type:          w.Type(),
			Configuration: config,
		}
		availableWidgets = append(availableWidgets, r)
		initialized[w.Identifier()] = w
	}

	// Register widgets.
	if err := env.Server.Register(availableWidgets); err != nil {
		return nil, err, true
	}

	return initialized, nil, true
}

func update(env *environment, initialized, started map[string]widgets.Widget) error {
	requested, err := env.Server.GetRequestedWidgets()
	if err != nil {
		return err
	}

	results := make([]widgets.BulkElement, 0, len(requested.Widgets))

	for _, widgetIdentifier := range requested.Widgets {
		widget, found := started[widgetIdentifier]
		if !found {
			widget, found = initialized[widgetIdentifier]
			if !found {
				log.Println("Unknown requested widget identifier: " + widgetIdentifier)
				widget = nil
			} else if err := widget.Start(); err != nil {
				log.Println("Could not start widget", widgetIdentifier, ":", err)
			}
		}

		if widget != nil {
			started[widgetIdentifier] = widget

			if err := widget.Update(); err != nil {
				log.Printf("Can't update %v: %v", widgetIdentifier, err)
				widget = nil
			} else if widget.HasData() == false {
				widget = nil
			}
		}

		widgetRaw, err := json.Marshal(&widget)
		if err != nil {
			return fmt.Errorf("Can't marshal widget %s: %v", widget.Identifier(), err)
		}

		e := widgets.BulkElement{
			Identifier: widgetIdentifier,
			Widget:     widgetRaw,
		}
		if widget != nil {
			e.Type = widget.Type()
		}

		results = append(results, e)
	}

	if err := env.Server.PostWidgetBulkUpdate(results); err != nil {
		return err
	}

	return nil
}

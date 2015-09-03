package server

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hverr/status-dashboard/server/settings"
)

type ClientConfiguration struct {
	Secret string `json:"secret"`
}

// Configuration holds the server configuration
var Configuration struct {
	UpdateInterval int                            `json:"updateInterval"`
	Clients        map[string]ClientConfiguration `json:"clients"`
	Users          map[string]string              `json:"users"`
}

// Validate a configuration if it is invalid an error is returned.
func ValidateConfiguration() error {
	return nil
}

// Parse options from a configuration file and validate the configuration file.
//
// Returns an error if reading the configuration file failed or if the resulting
// configuration could not be Validated.
func ParseConfiguration(file string) error {
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	if err := decoder.Decode(&Configuration); err != nil {
		return err
	}

	if v := Configuration.UpdateInterval; v > 0 {
		settings.MinimumClientUpdateInterval = time.Duration(v) * time.Second
		settings.MaximumClientUpdateInterval = time.Duration(3*v/2) * time.Second
	}

	return ValidateConfiguration()
}

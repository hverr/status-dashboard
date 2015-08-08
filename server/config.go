package server

import (
	"encoding/json"
	"os"
)

// Configuration holds the server configuration
var Configuration struct {
	Clients        []string            `json:"clients"`
	DefaultWidgets map[string][]string `json:"defaultWidgets"`
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
	if err := decoder.Decode(Configuration); err != nil {
		return err
	}

	return ValidateConfiguration()
}

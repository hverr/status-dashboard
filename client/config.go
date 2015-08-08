package client

import (
	"encoding/json"
	"errors"
	"os"
)

// Configuration holds the client configuration.
var Configuration struct {
	API        string   `json:"api"`
	Identifier string   `json:"identifier"`
	Name       string   `json:"name"`
	Widgets    []string `json:"widgets"`
}

// Validate a configuration. If it is invalid an error is returned.
func ValidateConfiguration() error {
	if Configuration.API == "" {
		return errors.New("No API is specified.")
	} else if Configuration.Widgets == nil {
		return errors.New("No widgets are specified.")
	}
	return nil
}

// Parse options from a configuration file and validate the configuration.
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

	return ValidateConfiguration()
}

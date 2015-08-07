package client

import (
	"encoding/json"
	"errors"
	"os"
)

// Configuration holds the client configuration.
type Configuration struct {
	API        string   `json:"api"`
	Identifier string   `json:"identifier"`
	Name       string   `json:"name"`
	Widgets    []string `json:"widgets"`
}

// NewDefaultConfiguration creates a new configuration with default values. It
// might not be a valid configuration.
func NewDefaultConfiguration() *Configuration {
	return &Configuration{}
}

// Validate a configuration. If it is invalid an error is returned.
func (c *Configuration) Validate() error {
	if c.API == "" {
		return errors.New("No API is specified.")
	} else if c.Widgets == nil {
		return errors.New("No widgets are specified.")
	}
	return nil
}

// Parse options from a configuration file and validate the configuration.
//
// Returns an error if reading the configuration file failed or if the resulting
// configuration could not be Validated.
func (c *Configuration) Parse(file string) error {
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	if err := decoder.Decode(c); err != nil {
		return err
	}

	return c.Validate()
}

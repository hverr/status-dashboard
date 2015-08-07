package server

import (
	"encoding/json"
	"os"
)

// Configuration holds the server configuration
type Configuration struct {
	Clients        []string            `json:"clients"`
	DefaultWidgets map[string][]string `json:"defaultWidgets"`
}

// NewDefaultConfiguration creates a new configuration with default values. It
// might not be a valid configuration.
func NewDefaultConfiguration() *Configuration {
	return &Configuration{}
}

// Validate a configuration if it is invalid an error is returned.
func (c *Configuration) Validate() error {
	return nil
}

// Parse options from a configuration file and validate the configuration file.
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

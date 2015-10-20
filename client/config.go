package client

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/hverr/status-dashboard/widgets"
)

// Configuration holds the client configuration.
type Configuration struct {
	API        string                     `json:"api"`
	Identifier string                     `json:"identifier"`
	Name       string                     `json:"name"`
	Secret     string                     `json:"secret"`
	Widgets    map[string]json.RawMessage `json:"widgets"`
}

// Validate a configuration. If it is invalid an error is returned.
func (c *Configuration) ValidateConfiguration() error {
	if c.API == "" {
		return errors.New("No API is specified.")
	} else if c.Widgets == nil {
		return errors.New("No widgets are specified.")
	}

	for w, _ := range c.Widgets {
		if widgets.AllWidgets[w] == nil {
			return errors.New("Unsupported widget " + w)
		}
	}

	return nil
}

// Parse options from a configuration file and validate the configuration.
//
// Returns an error if reading the configuration file failed or if the resulting
// configuration could not be Validated.
func (c *Configuration) ParseConfiguration(fh io.Reader) error {
	decoder := json.NewDecoder(fh)
	if err := decoder.Decode(&c); err != nil {
		return err
	}

	return c.ValidateConfiguration()
}

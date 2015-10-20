package server

import (
	"encoding/json"
	"os"
	"time"
)

type ClientConfiguration struct {
	Secret string `json:"secret"`
}

// Configuration holds the server configuration
type Configuration struct {
	UpdateInterval int                            `json:"updateInterval"`
	Clients        map[string]ClientConfiguration `json:"clients"`
	Users          map[string]string              `json:"users"`

	MinimumClientUpdateInterval time.Duration `json:"-"`
	MaximumClientUpdateInterval time.Duration `json:"-"`
	MaximumWidgetAge            time.Duration `json:"-"`
}

// Validate a configuration if it is invalid an error is returned.
func (c *Configuration) ValidateConfiguration() error {
	return nil
}

// Parse options from a configuration file and validate the configuration file.
//
// Returns an error if reading the configuration file failed or if the resulting
// configuration could not be Validated.
func (c *Configuration) ParseConfiguration(file string) error {
	fh, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	if err := decoder.Decode(&c); err != nil {
		return err
	}

	// Set defaults
	c.MinimumClientUpdateInterval = 3 * time.Second
	c.MaximumClientUpdateInterval = 5 * time.Second
	c.MaximumWidgetAge = 30 * time.Second

	if v := c.UpdateInterval; v > 0 {
		c.MinimumClientUpdateInterval = time.Duration(v) * time.Second
		c.MaximumClientUpdateInterval = time.Duration(3*v/2) * time.Second
	}

	return c.ValidateConfiguration()
}

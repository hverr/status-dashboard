package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	TestClientIdentifier = "webserver"
	TestClientName       = "Web Server"
	TestClientSecret     = "webserversecret"
)

func validConfig() Configuration {
	return Configuration{
		UpdateInterval: 1,
		Clients: map[string]ClientConfiguration{
			TestClientIdentifier: {
				Secret: TestClientSecret,
			},
		},
		Users: map[string]string{},
	}
}

func TestRegisterClient(t *testing.T) {
	r := ClientRegistration{
		Name:       TestClientName,
		Identifier: TestClientIdentifier,
		AvailableWidgets: []WidgetRegistration{
			{"load", []byte("null")},
			{"network", []byte(`{"interface":"eth0"}`)},
		},
	}

	// Test all valid
	{
		serv := NewServer(validConfig())
		err := serv.RegisterClient(&r)

		assert.Nil(t, err, "Unexpected error: %v", err)

		registered := serv.AllRegisteredClients()
		foundRegistration := false
		for _, i := range registered {
			if i.Identifier == TestClientIdentifier {
				foundRegistration = true
			}
		}
		assert.True(t, foundRegistration, "Client was not in registered list.")

		c, _ := serv.GetClient(TestClientIdentifier)
		assert.NotNil(t, c, "Client was not in initialized clients.")
		foundLoad := false
		foundNetwork := false
		for _, w := range c.AvailableWidgets {
			if w.Identifier() == "load" {
				foundLoad = true
			} else if w.Identifier() == "network_eth0" {
				foundNetwork = true
			}
		}
		assert.True(t, foundLoad, "Can't find load widget.")
		assert.True(t, foundNetwork, "Can't find network widget.")
	}

	// Test invalid widget
	{
		ir := ClientRegistration{
			TestClientName, TestClientIdentifier, []WidgetRegistration{
				{"invalid type", []byte("null")},
			},
		}
		serv := NewServer(validConfig())
		err := serv.RegisterClient(&ir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unknown widget type")
	}

	// Test invalid configuration json
	{
		ir := ClientRegistration{
			TestClientName, TestClientIdentifier, []WidgetRegistration{
				{"network", []byte("invalid json")},
			},
		}
		serv := NewServer(validConfig())
		err := serv.RegisterClient(&ir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Could not configure")
	}
}

func TestGetClient(t *testing.T) {
	// Client not found
	{
		serv := NewServer(validConfig())
		c, ok := serv.GetClient("unknown-client")
		assert.Nil(t, c)
		assert.False(t, ok)
	}

	// Valid query and expired client
	{
		r := ClientRegistration{
			Name:       TestClientName,
			Identifier: TestClientIdentifier,
			AvailableWidgets: []WidgetRegistration{
				{"load", []byte("null")},
			},
		}

		config := validConfig()
		config.MaximumWidgetAge = 10 * time.Second

		serv := NewServer(config)
		err := serv.RegisterClient(&r)
		assert.Nil(t, err, "Unexpected error: %v", err)

		c, _ := serv.GetClient(TestClientIdentifier)
		assert.NotNil(t, c)

		// Valid query
		c.LastSeen = time.Now()
		c, ok := serv.GetClient(TestClientIdentifier)
		assert.NotNil(t, c)
		assert.True(t, ok)

		// Expired
		c.LastSeen = time.Now().Add(-5 * time.Minute)
		c, ok = serv.GetClient(TestClientIdentifier)
		assert.NotNil(t, c)
		assert.False(t, ok)
	}
}

package server

import (
	"log"

	"github.com/hverr/status-dashboard/widgets"
	"github.com/pmylund/go-cache"
)

type Client struct {
	Name             string   `json:"name" binding:"required"`
	Identifier       string   `json:"identifier" binding:"required"`
	AvailableWidgets []string `json:"availableWidgets" binding:"required"`

	widgets          *cache.Cache `json:"-"`
	requestedWidgets []string     `json:"-"`
}

var RegisteredClients = cache.New(cache.NoExpiration, cache.NoExpiration)

func RegisterClient(client *Client) {
	log.Printf("Registering client %v (%v)\n", client.Name, client.Identifier)

	client.widgets = cache.New(cache.NoExpiration, cache.NoExpiration)
	client.requestedWidgets = []string{}

	defaultWidgets := Configuration.DefaultWidgets[client.Identifier]
	if len(defaultWidgets) == 0 {
		log.Print("Warning: client %v (%v) has no default widgets\n", client.Name, client.Identifier)
	}

	for _, w := range defaultWidgets {
		found := false
		for _, a := range client.AvailableWidgets {
			if w == a {
				found = true
				break
			}
		}

		if found {
			log.Printf("Registering widget %v for %v (%v)\n", w, client.Name, client.Identifier)
			client.requestedWidgets = append(client.requestedWidgets, w)

		} else {
			log.Printf("Warning: widget %v for %v (%v) is not available\n", w, client.Name, client.Identifier)
		}
	}

	RegisteredClients.Set(client.Identifier, client, cache.DefaultExpiration)
}

func GetClient(identifier string) (*Client, bool) {
	o, ok := RegisteredClients.Get(identifier)
	if !ok {
		return nil, false
	}
	return o.(*Client), true
}

func (c *Client) SetWidget(w widgets.Widget) {
	c.widgets.Set(w.Type(), w, cache.DefaultExpiration)
}

func (c *Client) RequestedWidgets() []string {
	return c.requestedWidgets
}

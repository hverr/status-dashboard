package server

import (
	"errors"
	"log"
	"time"

	"github.com/hverr/status-dashboard/widgets"
	"github.com/pmylund/go-cache"
)

var UnknownClientError = errors.New("Client unknown to the server.")

type Client struct {
	Name             string    `json:"name" binding:"required"`
	Identifier       string    `json:"identifier" binding:"required"`
	AvailableWidgets []string  `json:"availableWidgets" binding:"required"`
	LastSeen         time.Time `json:"-"`

	widgets          *cache.Cache `json:"-"`
	requestedWidgets []string     `json:"-"`
}

var RegisteredClients = cache.New(cache.NoExpiration, cache.NoExpiration)

func RegisterClient(client *Client) error {
	client.widgets = cache.New(cache.NoExpiration, cache.NoExpiration)

	log.Printf("Client registration request %v (%v)", client.Name, client.Identifier)

	// Make sure the client is allowed to connect.
	found := false
	for _, c := range Configuration.Clients {
		if c == client.Identifier {
			found = true
			break
		}
	}

	if !found {
		return UnknownClientError
	}

	// Actually register the client.
	client.LastSeen = time.Now()
	RegisteredClients.Set(client.Identifier, client, cache.DefaultExpiration)

	return nil
}

func AllRegisteredClients() []*Client {
	items := RegisteredClients.Items()
	result := make([]*Client, 0, len(items))
	for _, item := range items {
		result = append(result, item.Object.(*Client))
	}

	return result
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

func (c *Client) GetWidget(widgetType string) widgets.Widget {
	o, ok := c.widgets.Get(widgetType)
	if !ok {
		return nil
	}

	return o.(widgets.Widget)
}

func (c *Client) RequestedWidgets() []string {
	return c.requestedWidgets
}

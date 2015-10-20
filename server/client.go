package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server/settings"
	"github.com/hverr/status-dashboard/widgets"
	"github.com/pmylund/go-cache"
)

type Client struct {
	Name             string
	Identifier       string
	AvailableWidgets []widgets.Widget
	LastSeen         time.Time

	widgets *cache.Cache
}

type ClientRegistration struct {
	Name             string               `json:"name" binding:"required"`
	Identifier       string               `json:"identifier" binding:"required"`
	AvailableWidgets []WidgetRegistration `json:"availableWidgets" binding:"required"`
}

type WidgetRegistration struct {
	Type          string          `json:"type" binding:"required"`
	Configuration json.RawMessage `json:"configuration" binding:"required"`
}

func (s *server) RegisterClient(r *ClientRegistration) error {
	client := Client{
		Name:       r.Name,
		Identifier: r.Identifier,
		LastSeen:   time.Now(),
		widgets:    cache.New(cache.NoExpiration, cache.NoExpiration),
	}
	for _, widgetReg := range r.AvailableWidgets {
		initiator := widgets.AllWidgets[widgetReg.Type]
		if initiator == nil {
			return fmt.Errorf("Unknown widget type: %s", widgetReg.Type)
		}

		w := initiator()
		if err := w.Configure(widgetReg.Configuration); err != nil {
			return fmt.Errorf("Could not configure %s: %v", widgetReg.Type, err)
		}

		client.AvailableWidgets = append(client.AvailableWidgets, w)
	}

	log.Printf("Client registration request %v (%v)", client.Name, client.Identifier)

	// Actually register the client.
	s.registeredClients.Set(r.Identifier, r, cache.DefaultExpiration)
	s.initializedClients.Set(client.Identifier, &client, cache.DefaultExpiration)

	return nil
}

func (s *server) AllRegisteredClients() []*ClientRegistration {
	items := s.registeredClients.Items()
	result := make([]*ClientRegistration, 0, len(items))
	for _, item := range items {
		result = append(result, item.Object.(*ClientRegistration))
	}

	return result
}

func (s *server) GetClient(identifier string) (*Client, bool) {
	o, ok := s.initializedClients.Get(identifier)
	if !ok {
		return nil, false
	}

	c := o.(*Client)
	if time.Since(c.LastSeen) > settings.MaximumWidgetAge {
		return c, false
	}

	return c, true
}

func (s *server) AuthenticateClient(c *gin.Context, clientIdentifier string) bool {
	clientSecret := c.Request.Header.Get("X-Client-Secret")
	clientConfig, ok := s.Configuration.Clients[clientIdentifier]
	if !ok || clientSecret != clientConfig.Secret {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Could not authenticate client.",
		})
		return false
	}

	return true
}

func (c *Client) SetWidget(w widgets.Widget) {
	c.LastSeen = time.Now()
	c.widgets.Set(w.Identifier(), w, cache.DefaultExpiration)
}

func (c *Client) DeleteWidget(widgetIdentifier string) {
	c.widgets.Delete(widgetIdentifier)
}

func (c *Client) GetWidget(widgetIdentifier string) widgets.Widget {
	o, ok := c.widgets.Get(widgetIdentifier)
	if !ok {
		return nil
	}

	return o.(widgets.Widget)
}

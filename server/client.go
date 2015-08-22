package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server/settings"
	"github.com/hverr/status-dashboard/widgets"
	"github.com/pmylund/go-cache"
)

type Client struct {
	Name             string    `json:"name" binding:"required"`
	Identifier       string    `json:"identifier" binding:"required"`
	AvailableWidgets []string  `json:"availableWidgets" binding:"required"`
	LastSeen         time.Time `json:"-"`

	widgets *cache.Cache `json:"-"`
}

var RegisteredClients = cache.New(cache.NoExpiration, cache.NoExpiration)

func RegisterClient(client *Client) {
	client.widgets = cache.New(cache.NoExpiration, cache.NoExpiration)

	log.Printf("Client registration request %v (%v)", client.Name, client.Identifier)

	// Actually register the client.
	client.LastSeen = time.Now()
	RegisteredClients.Set(client.Identifier, client, cache.DefaultExpiration)
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

	c := o.(*Client)
	if time.Since(c.LastSeen) > settings.MaximumWidgetAge {
		return c, false
	}

	return c, true
}

func AuthenticateClient(c *gin.Context, clientIdentifier string) bool {
	clientSecret := c.Request.Header.Get("X-Client-Secret")
	clientConfig, ok := Configuration.Clients[clientIdentifier]
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
	c.widgets.Set(w.Type(), w, cache.DefaultExpiration)
}

func (c *Client) DeleteWidget(widgetType string) {
	c.widgets.Delete(widgetType)
}

func (c *Client) GetWidget(widgetType string) widgets.Widget {
	o, ok := c.widgets.Get(widgetType)
	if !ok {
		return nil
	}

	return o.(widgets.Widget)
}

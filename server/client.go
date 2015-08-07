package server

import "github.com/pmylund/go-cache"

type Client struct {
	Name             string `json:"name" binding:"required"`
	Identifier       string `json:"identifier" binding:"required"`
	AvailableWidgets string `json:"availableWidgets" binding:"required"`
}

var RegisteredClients = cache.New(cache.NoExpiration, cache.NoExpiration)

func RegisterClient(client *Client) {
	RegisteredClients.Set(client.Identifier, client, cache.DefaultExpiration)
}

func GetClient(identifier string) (*Client, bool) {
	o, ok := RegisteredClients.Get(identifier)
	if !ok {
		return nil, false
	}
	return o.(*Client), true
}

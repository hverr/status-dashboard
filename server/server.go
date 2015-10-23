package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
)

type Server interface {
	RegisterClient(r *ClientRegistration) error
	AllRegisteredClients() []*ClientRegistration
	GetClient(identifier string) (*Client, bool)
	AuthenticateClient(c *gin.Context, clientIdentifier string) bool
}

type server struct {
	Configuration Configuration

	registeredClients  *cache.Cache
	initializedClients *cache.Cache
}

func NewServer(c Configuration) Server {
	return &server{
		Configuration: c,

		registeredClients:  cache.New(cache.NoExpiration, cache.NoExpiration),
		initializedClients: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

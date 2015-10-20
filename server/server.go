package server

import "github.com/pmylund/go-cache"

type Server struct {
	Configuration Configuration

	registeredClients  *cache.Cache
	initializedClients *cache.Cache
}

func NewServer(c Configuration) *Server {
	return &Server{
		Configuration: c,

		registeredClients:  cache.New(cache.NoExpiration, cache.NoExpiration),
		initializedClients: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

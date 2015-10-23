package scheduler

import (
	"github.com/hverr/status-dashboard/server"
	"github.com/pmylund/go-cache"
)

type Scheduler interface {
	// Methods used to manage clients.
	RegisterClient(client string)
	HasClient(client string) bool
	RequestUpdateRequest(client string) chan []string
	FulfillUpdateRequest(client string, updated []string)

	// Methods used by the web api
	RequestWidgets(client string, widgets []string, immediately bool) chan bool
}

type scheduler struct {
	configuration server.Configuration

	widgetRequests *cache.Cache
}

func New(c server.Configuration) Scheduler {
	return &scheduler{
		configuration:  c,
		widgetRequests: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

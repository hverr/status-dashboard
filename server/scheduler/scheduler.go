package scheduler

import "github.com/pmylund/go-cache"

type Scheduler interface {
	// Methods used to manage clients.
	RegisterClient(client string)
	RequestUpdateRequest(client string) chan []string
	FulfillUpdateRequest(client string, updated []string)

	// Methods used by the web api
	RequestWidgets(client string, widgets []string, immediately bool) chan bool
}

type scheduler struct {
	widgetRequests *cache.Cache
}

func New() Scheduler {
	return &scheduler{
		widgetRequests: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

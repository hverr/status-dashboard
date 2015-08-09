package scheduler

import (
	"time"

	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/settings"

	"github.com/pmylund/go-cache"
)

var scheduler = cache.New(cache.NoExpiration, cache.NoExpiration)

func UpdateIntervalForClient(client *server.Client) time.Duration {
	err := scheduler.Add(client.Identifier, client, cache.DefaultExpiration)
	if err == nil {
		// Item was not yet in cache
		return 0
	}

	return settings.ClientUpdateInterval
}

func RegisterClient(client *server.Client) {
	scheduler.Delete(client.Identifier)
}

package scheduler

import (
	"sync"
	"time"

	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/settings"

	"github.com/pmylund/go-cache"
)

var scheduler = cache.New(cache.NoExpiration, cache.NoExpiration)

var updateListenerDispatcher = sync.Once{}
var updateListenerIncoming = make(chan bool)
var updateListenersLock = sync.RWMutex{}
var updateListeners = make([]chan bool, 0)

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

func RegisterUpdateListener() chan bool {
	c := make(chan bool)

	updateListenersLock.Lock()
	updateListeners = append(updateListeners, c)
	updateListenersLock.Unlock()

	return c
}

func NotifyUpdateListeners() {
	updateListenerDispatcher.Do(func() {
		go updateListenerBroadcaster()
	})

	updateListenerIncoming <- true
}

func updateListenerBroadcaster() {
	for {
		if flag := <-updateListenerIncoming; flag == false {
			// Channel has been closed.
			return
		}

		updateListenersLock.Lock()
		for _, listener := range updateListeners {
			listener <- true
		}
		updateListeners = make([]chan bool, 0)
		updateListenersLock.Unlock()
	}
}

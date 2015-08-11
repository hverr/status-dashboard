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

var widgetRequestDispatcher = sync.Once{}
var widgetRequestIncoming = make(chan int)
var widgetRequestCounter = 0
var widgetRequestListenersLock = sync.RWMutex{}
var widgetRequestListeners = make([]chan bool, 0)

func RegisterClientUpdateListener(client *server.Client) chan bool {
	c := make(chan bool, 1)

	err := scheduler.Add(client.Identifier, client, cache.DefaultExpiration)
	if err == nil {
		// Item was not yet in cache
		c <- true
		return c
	}

	<-time.After(settings.MinimumClientUpdateInterval)

	go func() {
		select {
		case <-RegisterWidgetRequestListener():
			c <- true
		case <-time.After(settings.MaximumClientUpdateInterval):
			c <- true
		}
	}()

	return c
}

func RegisterClient(client *server.Client) {
	scheduler.Delete(client.Identifier)
}

func RegisterWidgetRequest() {
	widgetRequestDispatcher.Do(func() {
		go widgetRequestBroadcaster()
	})

	widgetRequestIncoming <- 1
}

func DeregisterWidgetRequest() {
	widgetRequestDispatcher.Do(func() {
		go widgetRequestBroadcaster()
	})

	widgetRequestIncoming <- -1
}

func RegisterWidgetRequestListener() chan bool {
	c := make(chan bool, 1)

	widgetRequestListenersLock.Lock()
	widgetRequestListeners = append(widgetRequestListeners, c)
	widgetRequestListenersLock.Unlock()

	widgetRequestIncoming <- 0

	return c
}

func RegisterUpdateListener() chan bool {
	c := make(chan bool, 1)

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

func widgetRequestBroadcaster() {
	for {
		update := <-widgetRequestIncoming
		widgetRequestCounter += update

		if widgetRequestCounter > 0 {
			widgetRequestListenersLock.Lock()
			for _, listener := range widgetRequestListeners {
				listener <- true
			}
			widgetRequestListeners = make([]chan bool, 0)
			widgetRequestListenersLock.Unlock()
		}
	}
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

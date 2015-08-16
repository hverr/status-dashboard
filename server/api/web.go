package api

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/broadcaster"
	"github.com/hverr/status-dashboard/server/scheduler"
	"github.com/hverr/status-dashboard/server/settings"
)

func availableClients(c *gin.Context) {
	c.JSON(200, server.AllRegisteredClients())
}

func updateRequest(c *gin.Context) {
	request := make(map[string][]string)

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if len(request) == 0 {
		<-time.After(1 * time.Minute)
		c.JSON(200, gin.H{})
		return
	}

	// Inform the scheduler about all requested widgets
	fulfilled := make([]chan bool, 0, len(request))
	for client, widgets := range request {
		r := scheduler.RequestWidgets(client, widgets)

		if r == nil {
			// The client is unavailable. Pretend querying the client
			// to not overload the system.
			r = make(chan bool, 1)
			go func() {
				<-time.After(settings.MinimumClientUpdateInterval)
				r <- true
			}()
		}

		fulfilled = append(fulfilled, r)
	}

	// Wait for all requests to be fulfilled. If that takes
	// more than double the minimum client update interval, break.
	var wg sync.WaitGroup
	stopper := broadcaster.New()
	for _, channel := range fulfilled {
		wg.Add(1)
		go func() {
			select {
			case <-channel:
			case <-stopper.Listen():
			}
			wg.Done()
		}()
	}

	go func() {
		<-time.After(2 * settings.MinimumClientUpdateInterval)
		stopper.Emit()
	}()

	wg.Wait()

	result := make(map[string]map[string]interface{})
	for clientIdentifier, requestedWidgets := range request {
		clientResult := make(map[string]interface{})

		client, active := server.GetClient(clientIdentifier)

		for _, widgetType := range requestedWidgets {
			if client == nil || !active {
				clientResult[widgetType] = nil
			} else {
				clientResult[widgetType] = client.GetWidget(widgetType)
			}

		}

		result[clientIdentifier] = clientResult
	}

	c.JSON(200, result)
}

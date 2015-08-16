package scheduler

import (
	"time"

	"github.com/hverr/status-dashboard/server/settings"
	"github.com/pmylund/go-cache"
)

func RegisterClient(client string) {
	c := newWidgetRequestsContainer()
	widgetRequests.Add(client, c, cache.DefaultExpiration)
}

func RequestUpdateRequest(client string) chan []string {
	o, ok := widgetRequests.Get(client)
	if !ok {
		return nil
	}

	c := o.(*widgetRequestsContainer)

	out := make(chan []string, 1)
	go func() {
		minAge := settings.MinimumClientUpdateInterval
		age := time.Since(c.lastUpdated)

		if age < minAge && !c.hasImmediateRequest() {
			<-time.After(minAge - age)
		}

		for {
			s := c.requestedWidgets()
			if len(s) > 0 {
				out <- s
				return
			}
			// TODO: Fix race condition
			// Might have become dirty here
			<-c.dirty.Listen()
		}
	}()

	return out
}

func FulfillUpdateRequest(client string, updated []string) {
	o, ok := widgetRequests.Get(client)
	if !ok {
		return
	}

	c := o.(*widgetRequestsContainer)
	c.fulfillRequests(updated)
}

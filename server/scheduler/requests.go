package scheduler

import (
	"sync"
	"time"

	"github.com/hverr/status-dashboard/server/broadcaster"
	"github.com/pmylund/go-cache"
)

var widgetRequests = cache.New(cache.NoExpiration, cache.NoExpiration)

func RequestWidgets(client string, widgets []string, immediately bool) chan bool {
	r := newWidgetRequest(widgets)
	r.immediately = immediately

	o, ok := widgetRequests.Get(client)
	if !ok {
		return nil
	}
	c := o.(*widgetRequestsContainer)

	c.requestsLock.Lock()
	c.requests = append(c.requests, r)
	c.requestsLock.Unlock()

	c.dirty.Emit()

	return r.fulfilled
}

type widgetRequest struct {
	widgets     []string
	fulfilled   chan bool
	immediately bool
}

type widgetRequestsContainer struct {
	requestsLock sync.RWMutex
	requests     []*widgetRequest

	lastUpdated time.Time
	dirty       *broadcaster.Broadcaster
}

func newWidgetRequest(widgets []string) *widgetRequest {
	return &widgetRequest{
		widgets:   widgets,
		fulfilled: make(chan bool, 1),
	}
}

func newWidgetRequestsContainer() *widgetRequestsContainer {
	return &widgetRequestsContainer{
		dirty: broadcaster.New(),
	}
}

func (r *widgetRequest) isFulfilledBy(widgets []string) bool {
	for _, needle := range r.widgets {
		found := false
		for _, hay := range widgets {
			if hay == needle {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (c *widgetRequestsContainer) hasImmediateRequest() bool {
	flag := false

	c.requestsLock.RLock()
	for _, r := range c.requests {
		if r.immediately {
			flag = true
			break
		}
	}
	c.requestsLock.RUnlock()

	return flag
}

func (c *widgetRequestsContainer) requestedWidgets() []string {
	m := make(map[string]bool)

	c.requestsLock.RLock()
	for _, r := range c.requests {
		for _, w := range r.widgets {
			m[w] = true
		}
	}
	c.requestsLock.RUnlock()

	s := make([]string, 0, len(m))
	for key, _ := range m {
		s = append(s, key)
	}
	return s
}

func (c *widgetRequestsContainer) fulfillRequests(updated []string) {
	newRequests := make([]*widgetRequest, 0)

	c.requestsLock.Lock()

	for _, r := range c.requests {
		if !r.isFulfilledBy(updated) {
			newRequests = append(newRequests, r)
		} else {
			r.fulfilled <- true
		}
	}

	c.requests = newRequests
	c.lastUpdated = time.Now()

	c.requestsLock.Unlock()
}

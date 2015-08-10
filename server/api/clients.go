package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/scheduler"
	"github.com/hverr/status-dashboard/widgets"
)

func registerClient(c *gin.Context) {
	var client server.Client

	if err := c.BindJSON(&client); err != nil {
		c.AbortWithError(400, err)
	}

	server.RegisterClient(&client)
	scheduler.RegisterClient(&client)

	c.JSON(200, gin.H{})
}

func updateClientWidget(c *gin.Context) {
	client, ok := server.GetClient(c.Param("client"))
	if !ok {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	initiator := widgets.AllWidgets[c.Param("widget")]
	if initiator == nil {
		c.AbortWithError(404, errors.New("Widget not found."))
		return
	}

	widget := initiator()
	if err := c.BindJSON(&widget); err != nil {
		c.AbortWithError(400, err)
		return
	}

	client.SetWidget(widget)

	scheduler.NotifyUpdateListeners()

	c.JSON(200, gin.H{"status": "OK"})
}

func bulkUpdateClient(c *gin.Context) {
	client, ok := server.GetClient(c.Param("client"))
	if !ok {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	updates := make([]widgets.BulkElement, 0)
	if err := c.BindJSON(&updates); err != nil {
		c.AbortWithError(400, err)
	}

	result := make([]widgets.Widget, 0, len(updates))
	for _, u := range updates {
		initiator := widgets.AllWidgets[u.Type]
		if initiator == nil {
			c.AbortWithError(404, errors.New("Widget not found: "+u.Type))
			return
		}

		widget := initiator()
		encoded, err := json.Marshal(u.Widget)
		if err != nil {
			c.AbortWithError(500, err)
		}

		if err := json.Unmarshal(encoded, &widget); err != nil {
			c.AbortWithError(400, err)
			return
		}

		result = append(result, widget)
	}

	for _, w := range result {
		client.SetWidget(w)
	}

	scheduler.NotifyUpdateListeners()

	c.JSON(200, gin.H{"status": "OK"})
}

func requestedClientWidgets(c *gin.Context) {
	client, ok := server.GetClient(c.Param("client"))
	if !ok {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	delay := scheduler.UpdateIntervalForClient(client)
	<-time.After(delay)

	c.JSON(200, gin.H{"widgets": client.RequestedWidgets()})
}

func clientWidget(c *gin.Context) {
	client, ok := server.GetClient(c.Param("client"))
	if !ok {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	widget := client.GetWidget(c.Param("widget"))
	if widget == nil {
		c.AbortWithError(404, errors.New("Widget not found"))
		return
	}

	c.JSON(200, widget)
}

func allWidgets(c *gin.Context) {
	var request map[string][]string

	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if c.Query("force") != "true" {
		select {
		case <-scheduler.RegisterUpdateListener():
		case <-time.After(1 * time.Minute):
			c.JSON(202, gin.H{"reason": "No updates."})
			return
		}
	}

	result := make(map[string]map[string]widgets.Widget)

	for clientIdentifier, widgetTypes := range request {
		client, ok := server.GetClient(clientIdentifier)
		if !ok {
			c.AbortWithError(404, errors.New("Client not found: "+clientIdentifier))
			return
		}

		for _, widgetType := range widgetTypes {
			widget := client.GetWidget(widgetType)
			if widget == nil {
				c.AbortWithError(404, errors.New("Widget "+widgetType+
					" for "+clientIdentifier+"not found"))
				return
			}

			if result[clientIdentifier] == nil {
				result[clientIdentifier] = make(map[string]widgets.Widget)
			}
			result[clientIdentifier][widgetType] = widget
		}
	}

	c.JSON(200, result)
}

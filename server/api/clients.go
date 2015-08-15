package api

import (
	"encoding/json"
	"errors"

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

func bulkUpdateClient(c *gin.Context) {
	client, _ := server.GetClient(c.Param("client"))
	if client == nil {
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
	client, _ := server.GetClient(c.Param("client"))
	if client == nil {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	<-scheduler.RegisterClientUpdateListener(client)

	c.JSON(200, gin.H{"widgets": client.RequestedWidgets()})
}

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
	scheduler.RegisterClient(client.Identifier)

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

	updated := make([]string, 0, len(updates))
	result := make([]widgets.Widget, 0, len(updates))
	missing := make([]string, 0)
	for _, u := range updates {
		initiator := widgets.AllWidgets[u.Type]
		if initiator == nil {
			c.AbortWithError(404, errors.New("Widget not found: "+u.Type))
			return
		}

		if u.Widget != nil {
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

		} else {
			missing = append(missing, u.Type)
		}

		updated = append(updated, u.Type)
	}

	for _, w := range result {
		client.SetWidget(w)
	}

	for _, w := range missing {
		client.DeleteWidget(w)
	}

	scheduler.FulfillUpdateRequest(client.Identifier, updated)

	c.JSON(200, gin.H{"status": "OK"})
}

func requestedClientWidgets(c *gin.Context) {
	client, _ := server.GetClient(c.Param("client"))
	if client == nil {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	requested := <-scheduler.RequestUpdateRequest(client.Identifier)

	c.JSON(200, gin.H{"widgets": requested})
}

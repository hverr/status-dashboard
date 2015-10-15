package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/scheduler"
	"github.com/hverr/status-dashboard/server/settings"
	"github.com/hverr/status-dashboard/widgets"
)

func registerClient(c *gin.Context) {
	var r server.ClientRegistration

	if err := c.BindJSON(&r); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if !server.AuthenticateClient(c, r.Identifier) {
		return
	}

	if err := server.RegisterClient(&r); err != nil {
		c.AbortWithError(400, err)
		return
	}

	scheduler.RegisterClient(r.Identifier)

	c.JSON(200, gin.H{})
}

func bulkUpdateClient(c *gin.Context) {
	client, _ := server.GetClient(c.Param("client"))
	if client == nil {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	if !server.AuthenticateClient(c, client.Identifier) {
		return
	}

	updates := make([]widgets.BulkElement, 0)
	if err := c.BindJSON(&updates); err != nil {
		c.AbortWithError(400, err)
		return
	}

	updated := make([]string, 0, len(updates))
	result := make([]widgets.Widget, 0, len(updates))
	missing := make([]string, 0)
	for _, u := range updates {
		if u.Widget != nil && u.Type == "" {
			initiator := widgets.AllWidgets[u.Type]
			if initiator == nil {
				c.AbortWithError(404, errors.New("Widget not found: "+u.Type))
				return
			}

			widget := initiator()
			encoded, err := json.Marshal(u.Widget)
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			if err := json.Unmarshal(encoded, &widget); err != nil {
				c.AbortWithError(400, err)
				return
			}

			result = append(result, widget)

		} else {
			missing = append(missing, u.Type)
		}

		updated = append(updated, u.Identifier)
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

	if !server.AuthenticateClient(c, client.Identifier) {
		return
	}

	select {
	case requested := <-scheduler.RequestUpdateRequest(client.Identifier):
		c.JSON(200, gin.H{"widgets": requested})
	case <-time.After(settings.MaximumClientUpdateInterval):
		c.JSON(200, gin.H{"widgets": []string{}})
	}
}

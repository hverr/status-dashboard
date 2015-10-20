package api

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/settings"
	"github.com/hverr/status-dashboard/widgets"
)

func (api *API) registerClient(c *gin.Context) {
	var r server.ClientRegistration

	if err := c.BindJSON(&r); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if !api.Server.AuthenticateClient(c, r.Identifier) {
		return
	}

	if err := api.Server.RegisterClient(&r); err != nil {
		c.AbortWithError(400, err)
		return
	}

	api.Scheduler.RegisterClient(r.Identifier)

	c.JSON(200, gin.H{})
}

func (api *API) bulkUpdateClient(c *gin.Context) {
	client, _ := api.Server.GetClient(c.Param("client"))
	if client == nil {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	if !api.Server.AuthenticateClient(c, client.Identifier) {
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
		if u.Widget != nil && u.Type != "" {
			initiator := widgets.AllWidgets[u.Type]
			if initiator == nil {
				c.AbortWithError(404, errors.New("Widget not found: "+u.Type))
				return
			}

			widget := initiator()
			if err := json.Unmarshal(u.Widget, &widget); err != nil {
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

	api.Scheduler.FulfillUpdateRequest(client.Identifier, updated)

	c.JSON(200, gin.H{"status": "OK"})
}

func (api *API) requestedClientWidgets(c *gin.Context) {
	client, _ := api.Server.GetClient(c.Param("client"))
	if client == nil {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	if !api.Server.AuthenticateClient(c, client.Identifier) {
		return
	}

	select {
	case requested := <-api.Scheduler.RequestUpdateRequest(client.Identifier):
		c.JSON(200, gin.H{"widgets": requested})
	case <-time.After(settings.MaximumClientUpdateInterval):
		c.JSON(200, gin.H{"widgets": []string{}})
	}
}

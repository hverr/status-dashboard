package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/server/scheduler"
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

	scheduler.RegisterWidgetRequest()
	defer scheduler.DeregisterWidgetRequest()

	<-scheduler.RegisterUpdateListener()

	result := make(map[string]map[string]interface{})
	for clientIdentifier, requestedWidgets := range request {
		clientResult := make(map[string]interface{})

		client, found := server.GetClient(clientIdentifier)

		for _, widgetType := range requestedWidgets {
			if !found {
				clientResult[widgetType] = nil
			} else {
				clientResult[widgetType] = client.GetWidget(widgetType)
			}

		}

		result[clientIdentifier] = clientResult
	}

	c.JSON(200, result)
}

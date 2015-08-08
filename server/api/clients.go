package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/widgets"
)

func registerClient(c *gin.Context) {
	var client server.Client

	if err := c.BindJSON(&client); err != nil {
		c.AbortWithError(400, err)
	}

	server.RegisterClient(&client)

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
	c.JSON(200, gin.H{"status": "OK"})
}

func requestedClientWidgets(c *gin.Context) {
	client, ok := server.GetClient(c.Param("client"))
	if !ok {
		c.AbortWithError(404, errors.New("Client not found."))
		return
	}

	c.JSON(200, gin.H{"widgets": client.RequestedWidgets()})
}

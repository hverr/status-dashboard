package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
	"github.com/hverr/status-dashboard/widgets"
)

func registerClient(c *gin.Context) {
	var client server.Client

	if err := c.BindJSON(&client); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	server.RegisterClient(&client)

	c.JSON(200, gin.H{})
}

func updateClientWidget(c *gin.Context) {
	client, ok := server.GetClient(c.Param("client"))
	if !ok {
		c.JSON(404, gin.H{"error": "Client not found."})
		return
	}

	initiator := widgets.AllWidgets[c.Param("widget")]
	if initiator == nil {
		c.JSON(404, gin.H{"error": "Widget not found."})
		return
	}

	widget := initiator()
	if err := c.BindJSON(&widget); err != nil {
		c.JSON(400, gin.H{
			"error": "Could not decode widget: " + err.Error(),
		})
		return
	}

	client.SetWidget(widget)
	c.JSON(200, gin.H{"status": "OK"})
}

func requestedClientWidgets(c *gin.Context) {
	c.JSON(200, gin.H{
		"widgets": server.Configuration.DefaultWidgets[c.Param("client")],
	})
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
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

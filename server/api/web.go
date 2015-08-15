package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
)

func availableClients(c *gin.Context) {
	c.JSON(200, server.AllRegisteredClients())
}

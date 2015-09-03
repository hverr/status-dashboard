package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuthForUser(c *gin.Context) bool {
	u, p, _ := c.Request.BasicAuth()
	if AuthenticateUser(u, p) {
		return true
	}

	c.Header("WWW-Authenticate", "Basic realm=\"Status Dashboard\"")
	c.String(http.StatusUnauthorized, "Login required.")
	c.Abort()

	return false
}

func AuthenticateUser(username, password string) bool {
	log.Println("Authenticating", username, ":", password)
	log.Println("Users:", Configuration.Users)
	if Configuration.Users == nil {
		return true
	}

	pwd, ok := Configuration.Users[username]
	if ok && pwd == password {
		return true
	}

	return false
}

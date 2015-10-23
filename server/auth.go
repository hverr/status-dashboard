package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAuthenticator struct {
	Configuration Configuration
}

func (auth *UserAuthenticator) BasicAuthForUser(c *gin.Context) bool {
	u, p, _ := c.Request.BasicAuth()
	if auth.AuthenticateUser(u, p) {
		return true
	}

	c.Header("WWW-Authenticate", "Basic realm=\"Status Dashboard\"")
	c.String(http.StatusUnauthorized, "Login required.")
	c.Abort()

	return false
}

func (auth *UserAuthenticator) AuthenticateUser(username, password string) bool {
	if auth.Configuration.Users == nil {
		return true
	}

	pwd, ok := auth.Configuration.Users[username]
	if ok && pwd == password {
		return true
	}

	return false
}

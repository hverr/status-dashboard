package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server"
)

type API struct {
	Server            *server.Server
	UserAuthenticator server.UserAuthenticator
}

// Install installs the API end points.
//
// Some end points are used by the angular app to query for widget information,
// others are used by clients to push widget information.
func (api *API) Install(engine *gin.Engine) error {
	// Client API
	engine.POST("/api/clients/:client/register", api.registerClient)
	engine.POST("/api/clients/:client/bulk_update", api.bulkUpdateClient)
	engine.GET("/api/clients/:client/requested_widgets", api.requestedClientWidgets)

	// Angular API
	engine.GET("/api/available_clients", api.availableClients)
	engine.POST("/api/update_request", api.updateRequest)

	return nil
}

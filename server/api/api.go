package api

import "github.com/gin-gonic/gin"

// Install installs the API end points.
//
// Some end points are used by the angular app to query for widget information,
// others are used by clients to push widget information.
func Install(engine *gin.Engine) error {
	// Client API
	engine.POST("/api/clients/:client/register", registerClient)
	engine.POST("/api/clients/:client/bulk_update", bulkUpdateClient)
	engine.GET("/api/clients/:client/requested_widgets", requestedClientWidgets)

	// Angular API
	engine.GET("/api/available_clients", availableClients)
	engine.POST("/api/update_request", updateRequest)

	return nil
}

package api

import "github.com/gin-gonic/gin"

// Install installs the API end points.
//
// Some end points are used by the angular app to query for widget information,
// others are used by clients to push widget information.
func Install(engine *gin.Engine) error {
	engine.POST("/api/clients/:client/register", registerClient)
	engine.POST("/api/clients/:client/widgets/:widget/update", updateClientWidget)
	engine.GET("/api/clients/:client/widgets/requested", requestedClientWidgets)

	return nil
}

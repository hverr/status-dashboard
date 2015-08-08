package api

import "github.com/gin-gonic/gin"

// Install installs the API end points.
//
// Some end points are used by the angular app to query for widget information,
// others are used by clients to push widget information.
func Install(engine *gin.Engine) error {

	// POST /api/clients/register -> registerClient
	// POST /api/clients/:client/widgets/:widget/update -> updateClientWidget
	engine.POST("/api/clients/:client/widgets/:widget/update", func(c *gin.Context) {
		if c.Request.RequestURI == "/api/clients/register" {
			registerClient(c)
		} else {
			updateClientWidget(c)
		}
	})
	engine.GET("/api/clients/:client/widgets/requested", requestedClientWidgets)

	return nil
}

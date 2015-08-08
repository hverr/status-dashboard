package static

import (
	"os"

	"github.com/hverr/status-dashboard/server/settings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func Install(engine *gin.Engine) error {
	if _, err := os.Stat(settings.StaticAppRoot); err != nil {
		return err
	}

	engine.Use(static.Serve("/", static.LocalFile(settings.StaticAppRoot, true)))

	return nil
}

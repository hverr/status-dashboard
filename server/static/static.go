package static

import (
	"log"
	"os"

	"github.com/hverr/status-dashboard/server/settings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func Install(engine *gin.Engine) error {
	var err error
	var root string

	fs, err := LoadAssetFileSystem("/dist", true)
	if err == nil {
		log.Println("Serving static content from binary")
		engine.Use(static.Serve("/", fs))

	} else {
		log.Println("warning: could not read assets from binary:", err)

		toTry := []string{
			settings.StaticAppRoot,
			"./dist",
		}
		if envRoot := os.Getenv("HTML_ROOT"); envRoot != "" {
			toTry = append([]string{envRoot}, toTry...)
		}

		for _, path := range toTry {
			if _, err = os.Stat(path); err != nil {
				log.Println("warning: could not serve from", path)
			} else {
				root = path
				break
			}
		}

		if err != nil {
			return err
		}

		log.Println("Serving static content from", root)
		engine.Use(static.Serve("/", static.LocalFile(root, true)))
	}

	return nil
}

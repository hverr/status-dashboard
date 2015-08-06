package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hverr/status-dashboard/server/settings"
	"github.com/hverr/status-dashboard/server/static"
)

func main() {
	router := gin.Default()

	if err := static.Install(router); err != nil {
		log.Println("fatal: could not serve static fiels:", err)
		os.Exit(1)
	}

	router.Run(settings.ListenAddress)
}

package main

import (
	"log"
	"os"

	"github.com/hverr/status-dashboard/settings"
	"github.com/hverr/status-dashboard/static"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	if err := static.Install(router); err != nil {
		log.Println("fatal: could not serve static fiels:", err)
		os.Exit(1)
	}

	router.Run(settings.ListenAddress)
}

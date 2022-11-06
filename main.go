package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"web-app/database"
	"web-app/features"
	"web-app/shared/config"
	"web-app/web"
)

func main() {
	config.Init()

	if parseArguments() {
		return
	}

	database.Open(config.GetDbConnectionString())

	app := gin.New()
	app.Use(gin.Logger(), web.CustomRecovery())

	features.ConfigureRoutes(app)

	app.Run(fmt.Sprintf(":%d", config.GetAppPort()))
}

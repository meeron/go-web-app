package main

import (
	"github.com/gin-gonic/gin"
	"web-app/database"
	"web-app/features"
	"web-app/web"
)

func main() {
	if parseArguments() {
		return
	}

	database.Open()

	app := gin.New()
	app.Use(gin.Logger(), web.CustomRecovery())

	features.ConfigureRoutes(app)

	app.Run()
}

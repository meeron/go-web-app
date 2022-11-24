package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"web-app/database"
	docs "web-app/docs"
	"web-app/features"
	"web-app/shared/config"
	"web-app/shared/logger"
	"web-app/web"
)

func main() {
	config.Init()
	loggerWriter := logger.Init()

	appLogger := logger.Create("App")

	if parseArguments() {
		return
	}

	appLogger.Info("Connecting to database...")
	database.Open(config.GetDbConnectionString())
	appLogger.Info("Connected")

	if config.IsEnv(config.EnvProd) {
		gin.DisableConsoleColor()
	}

	gin.DefaultWriter = loggerWriter

	gin.SetMode(config.GetGinMode())

	app := gin.New()
	app.Use(gin.Logger(), web.CustomRecovery())

	features.ConfigureRoutes(app)

	address := fmt.Sprintf(":%d", config.GetAppPort())

	appLogger.Info("Listening on %v...", address)

	docs.SwaggerInfo.BasePath = "/"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	app.Run(address)
}

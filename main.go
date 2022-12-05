package main

import (
	"fmt"
	//swaggerfiles "github.com/swaggo/files"
	"web-app/database"
	"web-app/features"
	"web-app/shared/config"
	"web-app/shared/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.Init()
	logger.Init()

	appLogger := logger.Create("App")

	if parseArguments() {
		return
	}

	appLogger.Info("Connecting to database...")
	database.Open(config.GetDbConnectionString())
	appLogger.Info("Connected")

	app := fiber.New()
	app.Use(recover.New())

	features.ConfigureRoutes(app)

	address := fmt.Sprintf(":%d", config.GetAppPort())

	appLogger.Info("Listening on %v...", address)

	//docs.SwaggerInfo.BasePath = "/"
	//app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	app.Listen(address)
}

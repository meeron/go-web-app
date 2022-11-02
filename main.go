package main

import (
	"github.com/gin-gonic/gin"
	"web-app/database"
	"web-app/products"
	"web-app/users"
)

func main() {
	if parseArguments() {
		return
	}

	database.Open()

	app := gin.Default()

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"version": "1.0.0",
		})
	})

	app.Use()

	products.CreateRoutes(app, Auth())
	users.CreateRoutes(app, Auth())

	app.Run()
}

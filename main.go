package main

import (
	"github.com/gin-gonic/gin"
	"web-app/products"
)

func main() {
	parseArguments()

	app := gin.Default()

	app.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"version": "1.0.0",
		})
	})

	app.Use()

	products.CreateRoutes(app, Auth())

	app.Run()
}

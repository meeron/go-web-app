package home

import "github.com/gin-gonic/gin"

func index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": "1.0.0",
	})
}

func ConfigureRoutes(app *gin.Engine) {
	app.GET("/", index)
}

package home

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"web-app/shared"
)

func index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": "1.0.1",
	})
}

// Test go performance
// @Summary Tests go performance
// @Schemes
// @Description Do test
// @Tags Home
// @Produce json
// @Success 200
// @Param count path int true "Number of iteration"
// @Router /test/{count} [get]
func test(ctx *gin.Context) {
	cnt, err := strconv.Atoi(ctx.Param("count"))
	shared.PanicOnErr(err)

	count := 0

	for i := 0; i < cnt; i++ {
		count++
	}

	ctx.JSON(200, gin.H{
		"count": count,
	})
}

func ConfigureRoutes(app *gin.Engine) {
	app.GET("/", index)
	app.GET("/test/:count", test)
}

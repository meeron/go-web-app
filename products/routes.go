package products

import "github.com/gin-gonic/gin"

func Routes(app *gin.Engine) {
	app.GET("/products", GetAll)
	app.POST("/products", Add)
	app.GET("/products/:id", Get)
	app.DELETE("/products/:id", Delete)
}

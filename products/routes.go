package products

import "github.com/gin-gonic/gin"

func CreateRoutes(app *gin.Engine, middleware ...gin.HandlerFunc) *gin.RouterGroup {
	products := app.Group("/products")

	products.Use(middleware...)
	{
		products.GET("", GetAll)
		products.POST("", Add)
		products.GET(":id", Get)
		products.DELETE(":id", Delete)
	}

	return products
}

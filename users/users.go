package users

import (
	"github.com/gin-gonic/gin"
	"web-app/users/endpoints"
)

func CreateRoutes(app *gin.Engine, middleware ...gin.HandlerFunc) {
	app.POST("/login", endpoints.Login)

	usersGroup := app.Group("/users", middleware...)
	{
		usersGroup.POST("", endpoints.Create)
	}
}

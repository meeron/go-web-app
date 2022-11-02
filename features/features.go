package features

import (
	"github.com/gin-gonic/gin"
	"web-app/features/home"
	"web-app/features/products"
	"web-app/features/users"
)

func ConfigureRoutes(app *gin.Engine) {
	home.ConfigureRoutes(app)
	users.ConfigureRoutes(app)
	products.ConfigureRoutes(app)
}

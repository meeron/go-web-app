package features

import (
	"web-app/features/home"
	"web-app/features/products"
	"web-app/features/users"

	"github.com/gofiber/fiber/v2"
)

func ConfigureRoutes(app *fiber.App) {
	home.ConfigureRoutes(app)
	users.ConfigureRoutes(app)
	products.ConfigureRoutes(app)
}

package features

import (
	"web-app/features/home"

	"github.com/gofiber/fiber/v2"
)

func ConfigureRoutes(app *fiber.App) {
	home.ConfigureRoutes(app)
	//users.ConfigureRoutes(app)
	//products.ConfigureRoutes(app)
}

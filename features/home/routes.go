package home

import "github.com/gofiber/fiber/v2"

func index(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"version": "1.0.1",
	})
}

func ConfigureRoutes(app *fiber.App) {
	app.Get("/", index)
}

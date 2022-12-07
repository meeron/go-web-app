package web

import (
	"net/http"
	"strings"
	"web-app/web/jwt"

	"github.com/gofiber/fiber/v2"
)

func Auth() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		authHeader, ok := headers["Authorization"]
		if !ok {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		schema := parts[0]
		token := parts[1]

		if schema != "Bearer" {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		if !jwt.Validate(token) {
			return ctx.SendStatus(http.StatusUnauthorized)
		}

		return ctx.Next()
	}
}

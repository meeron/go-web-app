package web

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"web-app/web/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	const route = "/"

	app := fiber.New()
	app.Use(Auth())
	app.Get(route, func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	t.Run("should response 401 when authorization header is not provided", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest("GET", route, nil)

		// Act
		res, _ := app.Test(req)

		// Assert
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("should response 401 when authorization header is invalid", func(t *testing.T) {
		// Arrange
		const authHeader = "test"

		req := httptest.NewRequest("GET", route, nil)
		req.Header.Set("Authorization", authHeader)

		// Act
		res, _ := app.Test(req)
		resAuthHeader := res.Request.Header.Get("Authorization")

		// Assert
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, authHeader, resAuthHeader)
	})

	t.Run("should response 401 when authorization header has wrong schema", func(t *testing.T) {
		// Arrange
		const authHeader = "test test"

		req := httptest.NewRequest("GET", route, nil)
		req.Header.Set("Authorization", authHeader)

		// Act
		res, _ := app.Test(req)
		resAuthHeader := res.Request.Header.Get("Authorization")

		// Assert
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, authHeader, resAuthHeader)
	})

	t.Run("should response 401 when authorization header has wrong token", func(t *testing.T) {
		// Arrange
		const authHeader = "Bearer test"

		req := httptest.NewRequest("GET", route, nil)
		req.Header.Set("Authorization", authHeader)

		// Act
		res, _ := app.Test(req)
		resAuthHeader := res.Request.Header.Get("Authorization")

		// Assert
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		assert.Equal(t, authHeader, resAuthHeader)
	})

	t.Run("should response 200 when authorization header is valid", func(t *testing.T) {
		// Arrange
		token := jwt.Create(make(map[string]string, 0))
		authHeader := fmt.Sprintf("Bearer %s", token)

		req := httptest.NewRequest("GET", route, nil)
		req.Header.Set("Authorization", authHeader)

		// Act
		res, _ := app.Test(req)

		// Assert
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}

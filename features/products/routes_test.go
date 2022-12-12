package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"web-app/database"
	"web-app/shared"
	"web-app/web/jwt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

const (
	querySelectAll = `SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL`
)

func TestGetAll(t *testing.T) {
	dbMock, db, _ := database.OpenMock()
	defer db.Close()

	app := fiber.New()
	ConfigureRoutes(app)

	t.Run("should response with 200 and list of products", func(t *testing.T) {
		// Arrange
		req := newRequest("GET", "/products", nil)

		rows := sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(1, "1", 1.0).
			AddRow(2, "2", 2.0)

		dbMock.ExpectQuery(regexp.QuoteMeta(querySelectAll)).
			WillReturnRows(rows)

		// Act
		res, _ := app.Test(req)
		products := make([]Product, 0)
		json.Unmarshal(shared.Unwrap(io.ReadAll(res.Body)), &products)

		// Assert
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Equal(t, 2, len(products))
		assert.Equal(t, uint(1), products[0].Id)
		assert.Equal(t, uint(2), products[1].Id)
		assert.Equal(t, "1", products[0].Name)
		assert.Equal(t, "2", products[1].Name)
		assert.Equal(t, float32(1.0), products[0].Price)
		assert.Equal(t, float32(2.0), products[1].Price)
	})
}

func newRequest(method string, url string, body interface{}) *http.Request {
	var bodyReader io.Reader = nil
	token := jwt.Create(make(map[string]string, 0))
	authHeader := fmt.Sprintf("Bearer %s", token)

	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	return req
}

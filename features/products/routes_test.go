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
	"web-app/tests"
	"web-app/web"
	"web-app/web/jwt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

const (
	querySelectAll = `SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL`
	insertQuery    = `INSERT INTO "products" ("created_at","updated_at","deleted_at","name","price") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
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

func TestAdd(t *testing.T) {
	dbMock, db, _ := database.OpenMock()
	defer db.Close()
	app := fiber.New()
	ConfigureRoutes(app)

	t.Run("should response with 400 when request body is invalid", func(t *testing.T) {
		// Arrange
		req := newRequest("POST", "/products", nil)

		// Act
		res, _ := app.Test(req)
		resErr := parseErr(res)

		// Assert
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "invalid request", resErr.Message)
	})

	t.Run("should response with 400 when name is invalid", func(t *testing.T) {
		// Arrange
		req := newRequest("POST", "/products", NewProduct{})

		// Act
		res, _ := app.Test(req)
		resErr := parseErr(res)

		// Assert
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "'Name' is required", resErr.Message)
	})

	t.Run("should response with 400 when price is invalid", func(t *testing.T) {
		// Arrange
		req := newRequest("POST", "/products", NewProduct{Name: "test"})

		// Act
		res, _ := app.Test(req)
		resErr := parseErr(res)

		// Assert
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "'Price' must be greater than 0", resErr.Message)
	})

	t.Run("should response with 200 request body is valid", func(t *testing.T) {
		// Arrange
		const name = "test"
		const price float32 = 6.66
		const addedId uint = 666

		body := NewProduct{Name: name, Price: price}

		req := newRequest("POST", "/products", body)

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(addedId)

		dbMock.ExpectBegin()
		dbMock.ExpectQuery(regexp.QuoteMeta(insertQuery)).
			WithArgs(
				tests.AnyTime{},
				tests.AnyTime{},
				nil,
				name,
				price,
			).WillReturnRows(rows)
		dbMock.ExpectCommit()

		// Act
		res, _ := app.Test(req)
		product := Product{}
		json.Unmarshal(shared.Unwrap(io.ReadAll(res.Body)), &product)

		// Assert
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Equal(t, addedId, product.Id)
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

func parseErr(res *http.Response) web.Error {
	resErr := web.Error{}
	json.Unmarshal(shared.Unwrap(io.ReadAll(res.Body)), &resErr)

	return resErr
}

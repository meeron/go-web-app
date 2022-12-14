package users

import (
	"bytes"
	"crypto/sha256"
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
	query       = `SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`
	queryExists = `SELECT count(*) FROM "users" WHERE email = $1`
	queryInsert = `INSERT INTO "users" ("created_at","updated_at","deleted_at","email","password") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
)

func TestLogin(t *testing.T) {
	app := fiber.New()
	ConfigureRoutes(app)

	dbMock, db, _ := database.OpenMock()
	defer db.Close()

	t.Run("should response with 400 when request is invalid", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest("POST", "/login", nil)
		req.Header.Set("Content-Type", "application/json")

		// Act
		res, _ := app.Test(req)
		resErr := web.Error{}
		json.Unmarshal(shared.Unwrap(io.ReadAll(res.Body)), &resErr)

		// Assert
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "invalid request", resErr.Message)
	})

	t.Run("should respone with 401 when user is not found", func(t *testing.T) {
		// Arrange
		const email = "test"

		login := Login{Email: email}
		jsonBytes, _ := json.Marshal(login)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBytes))
		req.Header.Set("Content-Type", "application/json")

		dbMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows(nil))

		// Act
		res, _ := app.Test(req)
		resErr := web.Error{}
		json.Unmarshal(shared.Unwrap(io.ReadAll(res.Body)), &resErr)

		// Assert
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("should respone with 200 when password is valid", func(t *testing.T) {
		// Arrange
		const email = "test"
		const password = "pass"

		login := Login{Email: email, Password: password}
		jsonBytes, _ := json.Marshal(login)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBytes))
		req.Header.Set("Content-Type", "application/json")

		hash := sha256.New()
		hash.Write([]byte(password))

		hashedPass := fmt.Sprintf("%x", hash.Sum(nil))

		rows := sqlmock.NewRows([]string{"password"}).
			AddRow(hashedPass)

		dbMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(email).
			WillReturnRows(rows)

		// Act
		res, _ := app.Test(req)

		// Assert
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestCreate(t *testing.T) {
	dbMock, db, _ := database.OpenMock()
	defer db.Close()

	app := fiber.New()
	ConfigureRoutes(app)

	t.Run("should response with 400 when request body is invalid", func(t *testing.T) {
		// Arrange
		body := NewUser{}

		req := newCreateRequest(body)

		// Act
		res, _ := app.Test(req)
		resErr := parseErr(res)

		// Assert
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "'Email' is required", resErr.Message)
	})

	t.Run("should response with 400 when email invalid", func(t *testing.T) {
		// Arrange
		body := NewUser{
			Email:    "invalid_email",
			Password: "test",
		}

		req := newCreateRequest(body)

		// Act
		res, _ := app.Test(req)
		resErr := parseErr(res)

		// Assert
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "'Email' is invalid", resErr.Message)
	})

	t.Run("should response with 422 when user exists", func(t *testing.T) {
		// Arrange
		const email = "test@test.com"

		body := NewUser{
			Email:    email,
			Password: "test",
		}

		req := newCreateRequest(body)

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)
		dbMock.ExpectQuery(regexp.QuoteMeta(queryExists)).
			WithArgs(email).
			WillReturnRows(rows)

		// Act
		res, _ := app.Test(req)
		resErr := parseErr(res)

		// Assert
		assert.Equal(t, fiber.StatusUnprocessableEntity, res.StatusCode)
		assert.Equal(t, "Exists", resErr.ErrorCode)
	})

	t.Run("should response with 200 when user has been added", func(t *testing.T) {
		// Arrange
		const email = "test_added@test.com"
		const password = "test"
		const addedId uint = 666

		hashedPass, _ := shared.Sha256(password)

		body := NewUser{
			Email:    email,
			Password: password,
		}

		req := newCreateRequest(body)

		dbMock.ExpectQuery(regexp.QuoteMeta(queryExists)).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows(nil))

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(addedId)

		dbMock.ExpectBegin()
		dbMock.ExpectQuery(regexp.QuoteMeta(queryInsert)).
			WithArgs(tests.AnyTime{}, tests.AnyTime{}, nil, email, hashedPass).
			WillReturnRows(rows)
		dbMock.ExpectCommit()

		// Act
		res, _ := app.Test(req)
		resUser := User{}
		json.Unmarshal(shared.Unwrap(io.ReadAll(res.Body)), &resUser)

		// Assert
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		assert.Equal(t, int(addedId), resUser.Id)
	})
}

func newCreateRequest(body NewUser) *http.Request {
	token := jwt.Create(make(map[string]string, 0))
	authHeader := fmt.Sprintf("Bearer %s", token)

	bodyJsonBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(bodyJsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	return req
}

func parseErr(res *http.Response) web.Error {
	resErr := web.Error{}
	json.Unmarshal(shared.Unwrap(io.ReadAll(res.Body)), &resErr)

	return resErr
}

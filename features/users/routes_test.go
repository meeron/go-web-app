package users

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"web-app/database"
	"web-app/tests"
	"web-app/web"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	query = `SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`
)

func TestLogin(t *testing.T) {
	dbMock, db, _ := database.OpenMock()
	defer db.Close()
	r := tests.SetUpRouter()
	r.POST("/", login)

	t.Run("should response with 400 when request is invalid", func(t *testing.T) {
		// Arrange
		req, _ := http.NewRequest("POST", "/", nil)
		res := httptest.NewRecorder()
		resErr := web.Error{}

		// Act
		r.ServeHTTP(res, req)
		json.Unmarshal(res.Body.Bytes(), &resErr)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "invalid request", resErr.Message)
	})

	t.Run("should respone with 401 when user is not found", func(t *testing.T) {
		// Arrange
		const email = "test1"
		login := Login{Email: email}
		jsonBytes, _ := json.Marshal(login)

		req, _ := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
		res := httptest.NewRecorder()

		dbMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows(nil))

		// Act
		r.ServeHTTP(res, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("should respone with 401 when password is invalid", func(t *testing.T) {
		// Arrange
		const email = "test"
		login := Login{Email: email, Password: "pass"}
		jsonBytes, _ := json.Marshal(login)

		rows := sqlmock.NewRows([]string{"password"}).
			AddRow("test_pass")

		req, _ := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
		res := httptest.NewRecorder()

		dbMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(email).
			WillReturnRows(rows)

		// Act
		r.ServeHTTP(res, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("should respone with 200 when password is valid", func(t *testing.T) {
		// Arrange
		const email = "test"
		const password = "pass"
		login := Login{Email: email, Password: password}
		jsonBytes, _ := json.Marshal(login)

		hash := sha256.New()
		hash.Write([]byte(password))

		hashedPass := fmt.Sprintf("%x", hash.Sum(nil))

		rows := sqlmock.NewRows([]string{"password"}).
			AddRow(hashedPass)

		req, _ := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
		res := httptest.NewRecorder()

		dbMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(email).
			WillReturnRows(rows)

		// Act
		r.ServeHTTP(res, req)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)
	})

	/*
		t.Run("should respone with proper access token", func(t *testing.T) {
			// Arrange
			const email = "test@test.pl"
			const password = "pass"
			const id uint = 666
			login := Login{Email: email, Password: password}
			jsonBytes, _ := json.Marshal(login)

			hash := sha256.New()
			hash.Write([]byte(password))

			hashedPass := fmt.Sprintf("%x", hash.Sum(nil))

			rows := sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(id, email, hashedPass)

			req, _ := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
			res := httptest.NewRecorder()

			dbMock.ExpectQuery(regexp.QuoteMeta(query)).
				WithArgs(email).
				WillReturnRows(rows)

			// Act
			r.ServeHTTP(res, req)
			token := Token{}
			json.Unmarshal(res.Body.Bytes(), &token)

			// Assert
			assert.Equal(t, http.StatusOK, res.Code)
			assert.NotEmpty(t, token.AccessToken)
		})
	*/
}

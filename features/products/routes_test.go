package products

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"web-app/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	getQuery = `SELECT * FROM "products" WHERE "products"."id" = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`
)

func TestGet(t *testing.T) {
	dbMock, _ := database.OpenMock()

	r := SetUpRouter()
	r.GET("/:id", get)

	t.Run("should response with 422 if no product found", func(t *testing.T) {
		dbMock.ExpectQuery(regexp.QuoteMeta(getQuery)).
			WithArgs(123).
			WillReturnRows(sqlmock.NewRows(nil))

		req, _ := http.NewRequest("GET", "/123", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})

	t.Run("should response 200", func(t *testing.T) {
		dbMock.ExpectQuery(regexp.QuoteMeta(getQuery)).
			WithArgs(666).
			WillReturnRows(sqlmock.NewRows(nil))

		req, _ := http.NewRequest("GET", "/666", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

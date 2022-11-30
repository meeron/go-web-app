package products

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"web-app/database"
	"web-app/web"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	getQuery = `SELECT * FROM "products" WHERE "products"."id" = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`
)

func TestGetProductById(t *testing.T) {
	dbMock, _ := database.OpenMock()

	r := SetUpRouter()
	r.GET("/:id", get)

	t.Run("should response with 422 when no product found", func(t *testing.T) {
		// Arrange
		dbMock.ExpectQuery(regexp.QuoteMeta(getQuery)).
			WithArgs(123).
			WillReturnRows(sqlmock.NewRows(nil))
		resErr := web.Error{}

		// Act
		req, _ := http.NewRequest("GET", "/123", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		json.Unmarshal(res.Body.Bytes(), &resErr)

		// Assert
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		assert.Equal(t, "NotFound", resErr.ErrorCode)
	})

	t.Run("should response 200 when product is found", func(t *testing.T) {
		// Arrange
		const id uint = 666
		const name = "Test product"
		const price float32 = 9.99

		expectedRow := sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(id, name, price)

		dbMock.ExpectQuery(regexp.QuoteMeta(getQuery)).
			WithArgs(id).
			WillReturnRows(expectedRow)
		product := Product{}

		// Act
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%d", id), nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		json.Unmarshal(res.Body.Bytes(), &product)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, id, product.Id)
		assert.Equal(t, name, product.Name)
		assert.Equal(t, price, product.Price)
	})
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

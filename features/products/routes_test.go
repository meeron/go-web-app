package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"web-app/database"
	"web-app/tests"
	"web-app/web"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	getByIdQuery = `SELECT * FROM "products" WHERE "products"."id" = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`
	getAllQuery  = `SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL`
	insertQuery  = `INSERT INTO "products" ("created_at","updated_at","deleted_at","name","price") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`
)

func TestGetProductById(t *testing.T) {
	dbMock, _ := database.OpenMock()

	r := SetUpRouter()
	r.GET("/:id", get)

	t.Run("should response with 422 when no product found", func(t *testing.T) {
		// Arrange
		const id uint = 123

		dbMock.ExpectQuery(regexp.QuoteMeta(getByIdQuery)).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(nil))
		resErr := web.Error{}

		// Act
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%d", id), nil)
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
		product := Product{}

		expectedRow := sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(id, name, price)

		dbMock.ExpectQuery(regexp.QuoteMeta(getByIdQuery)).
			WithArgs(id).
			WillReturnRows(expectedRow)
		req, _ := http.NewRequest("GET", fmt.Sprintf("/%d", id), nil)
		res := httptest.NewRecorder()

		// Act
		r.ServeHTTP(res, req)
		json.Unmarshal(res.Body.Bytes(), &product)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, id, product.Id)
		assert.Equal(t, name, product.Name)
		assert.Equal(t, price, product.Price)
	})

	t.Run("should response with 404 if id is not int", func(t *testing.T) {
		// Arrange
		req, _ := http.NewRequest("GET", "/abc", nil)
		res := httptest.NewRecorder()

		// Act
		r.ServeHTTP(res, req)

		// Assert
		assert.Equal(t, 404, res.Code)
	})
}

func TestGetAllProducts(t *testing.T) {
	dbMock, _ := database.OpenMock()

	r := SetUpRouter()
	r.GET("/", getAll)

	t.Run("should return empty list when there is no products", func(t *testing.T) {
		// Arrange
		products := make([]Product, 0)

		dbMock.ExpectQuery(regexp.QuoteMeta(getAllQuery)).
			WillReturnRows(sqlmock.NewRows(nil))

		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()

		// Act
		r.ServeHTTP(res, req)
		json.Unmarshal(res.Body.Bytes(), &products)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Empty(t, products)
	})

	t.Run("should return products when there some are found", func(t *testing.T) {
		// Arrange
		products := []Product{
			{Id: 1, Name: "test1", Price: 1.1},
			{Id: 2, Name: "test2", Price: 2.2},
		}

		rows := createRows().
			AddRow(products[0].Id, products[0].Name, products[0].Price).
			AddRow(products[1].Id, products[1].Name, products[1].Price)

		dbMock.ExpectQuery(regexp.QuoteMeta(getAllQuery)).
			WillReturnRows(rows)

		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()

		// Act
		r.ServeHTTP(res, req)

		result := make([]Product, 0)
		json.Unmarshal(res.Body.Bytes(), &result)

		// Assert
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Len(t, products, 2)
		assert.EqualValues(t, products, result)
	})
}

func TestAddProduct(t *testing.T) {
	dbMock, _ := database.OpenMock()

	r := SetUpRouter()
	r.POST("/", add)

	t.Run("should response with 400 when body is invalid", func(t *testing.T) {
		// Arrange
		req, _ := http.NewRequest("POST", "/", nil)
		res := httptest.NewRecorder()
		resErr := web.Error{}

		//Act
		r.ServeHTTP(res, req)
		json.Unmarshal(res.Body.Bytes(), &resErr)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "BadRequest", resErr.ErrorCode)
		assert.Equal(t, "Request body should be in JSON format", resErr.Message)
	})

	t.Run("should response with 500 when there add err", func(t *testing.T) {
		// Arrange
		body := strings.NewReader("{}")
		req, _ := http.NewRequest("POST", "/", body)
		res := httptest.NewRecorder()
		resErr := web.Error{}

		dbMock.ExpectBegin()
		dbMock.ExpectQuery(regexp.QuoteMeta(insertQuery)).
			WithArgs(tests.AnyTime{}, tests.AnyTime{}, nil, "", 0.0).
			WillReturnError(errors.New("not inserted"))
		dbMock.ExpectRollback()

		//Act
		r.ServeHTTP(res, req)
		json.Unmarshal(res.Body.Bytes(), &resErr)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	/*
		t.Run("should response with 201 when product is created", func(t *testing.T) {
			// Arrange
			const name = "new product"
			const price float32 = 9.99

			newProduct := NewProduct{
				Name:  "new_product",
				Price: price,
			}

			jsonBytes, _ := json.Marshal(newProduct)
			body := bytes.NewReader(jsonBytes)

			req, _ := http.NewRequest("POST", "/", body)
			res := httptest.NewRecorder()
			product := Product{}

			// dbMock.ExpectBegin()
			// dbMock.ExpectQuery("test")
			// dbMock.ExpectCommit()
			dbMock.ExpectExec("INSERT INTO products")

			//Act
			r.ServeHTTP(res, req)
			json.Unmarshal(res.Body.Bytes(), &product)

			// Assert
			assert.Equal(t, http.StatusCreated, res.Code)
			assert.Equal(t, name, product.Name)
			assert.Equal(t, price, product.Price)
		})
	*/
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	return router
}

func createRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "price"})
}

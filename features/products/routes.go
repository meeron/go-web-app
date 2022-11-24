package products

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web-app/database"
	"web-app/shared"
	"web-app/web"
)

// Get all products
// @Summary Get all products
// @Schemes
// @Tags Products
// @Produce json
// @Success 200 {array} products.Product
// @Router /products [get]
func getAll(ctx *gin.Context) {
	db := database.DbCtx()

	products := db.Products().Find()

	result := make([]Product, 0)

	for _, p := range products {
		result = append(result, Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

// @Summary Add product
// @Schemes
// @Tags Products
// @Produce json
// @Accept json
// @Param request body products.NewProduct true "New product"
// @Success 201 {object} products.Product
// @Router /products [post]
func add(ctx *gin.Context) {
	var body NewProduct

	shared.PanicOnErr(ctx.BindJSON(&body))

	db := database.DbCtx()

	newEntity := database.Product{
		Name:  body.Name,
		Price: body.Price,
	}

	db.Products().Add(&newEntity)

	ctx.JSON(http.StatusCreated, Product{
		Id:    newEntity.ID,
		Name:  newEntity.Name,
		Price: newEntity.Price,
	})
}

// @Summary Get product
// @Schemes
// @Tags Products
// @Produce json
// @Param id path int true "Product's id"
// @Success 200 {object} products.Product
// @Failure 422 {object} web.Error
// @Router /products/{id} [get]
func get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	shared.PanicOnErr(err)

	db := database.DbCtx()

	product := db.Products().GetById(id)
	if product == nil {
		ctx.JSON(http.StatusUnprocessableEntity, web.NotFound())
		return
	}

	ctx.JSON(http.StatusOK, Product{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	})
}

// @Summary Delete product
// @Schemes
// @Tags Products
// @Produce json
// @Param id path int true "Product's id"
// @Success 200 {object} products.Product
// @Failure 422 {object} web.Error
// @Router /products/{id} [delete]
func remove(ctx *gin.Context) {
	id := shared.Unwrap(strconv.Atoi(ctx.Param("id")))

	db := database.DbCtx()

	result := db.Products().Remove(id)
	if !result {
		ctx.JSON(422, web.NotFound())
		return
	}

	ctx.Status(http.StatusOK)
}

func ConfigureRoutes(app *gin.Engine) {
	g := app.Group("/products", web.Auth())
	{
		g.GET("", getAll)
		g.POST("", add)
		g.GET(":id", get)
		g.DELETE(":id", remove)
	}
}

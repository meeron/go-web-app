package products

import (
	"web-app/database"
)

// Get all products
// @Summary Get all products
// @Schemes
// @Tags Products
// @Produce json
// @Success 200 {array} products.Product
// @Router /products [get]
func getAll() {
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

	//ctx.JSON(http.StatusOK, result)
}

// @Summary Add product
// @Schemes
// @Tags Products
// @Produce json
// @Accept json
// @Param request body products.NewProduct true "New product"
// @Success 201 {object} products.Product
// @Router /products [post]
func add() {
	/*
		var body NewProduct

		bindErr := ctx.ShouldBindJSON(&body)
		if bindErr != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, web.BadRequest(bindErr))
			return
		}

		db := database.DbCtx()

		newEntity := database.Product{
			Name:  body.Name,
			Price: body.Price,
		}

		addErr := db.Products().Add(&newEntity)
		if addErr != nil {
			// TODO: Log error
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusCreated, Product{
			Id:    newEntity.ID,
			Name:  newEntity.Name,
			Price: newEntity.Price,
		})
	*/
}

// @Summary Get product
// @Schemes
// @Tags Products
// @Produce json
// @Param id path int true "Product's id"
// @Success 200 {object} products.Product
// @Failure 422 {object} web.Error
// @Router /products/{id} [get]
func get() {
	/*
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		db := database.DbCtx()

		product := db.Products().GetById(id)
		if product == nil {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.NotFound())
			return
		}

		ctx.JSON(http.StatusOK, Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		})
	*/
}

// @Summary Delete product
// @Schemes
// @Tags Products
// @Produce json
// @Param id path int true "Product's id"
// @Success 200 {object} products.Product
// @Failure 422 {object} web.Error
// @Router /products/{id} [delete]
func remove() {
	/*
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		db := database.DbCtx()

		result := db.Products().Remove(id)
		if !result {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, web.NotFound())
			return
		}

		ctx.Status(http.StatusOK)
	*/
}

func ConfigureRoutes() {
	/*
		g := app.Group("/products", web.Auth())
		{
			g.GET("", getAll)
			g.POST("", add)
			g.GET(":id", get)
			g.DELETE(":id", remove)
		}*/
}

package products

import (
	"log"
	"web-app/database"
	"web-app/shared"
	"web-app/web"

	"github.com/gofiber/fiber/v2"
)

// Get all products
// @Summary Get all products
// @Schemes
// @Tags Products
// @Produce json
// @Success 200 {array} products.Product
// @Router /products [get]
func getAll(ctx *fiber.Ctx) error {
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

	return ctx.JSON(result)
}

// @Summary Add product
// @Schemes
// @Tags Products
// @Produce json
// @Accept json
// @Param request body products.NewProduct true "New product"
// @Success 201 {object} products.Product
// @Router /products [post]
func add(ctx *fiber.Ctx) error {
	var body NewProduct

	if bindErr := ctx.BodyParser(&body); bindErr != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(web.BadRequest(bindErr))
	}

	if valErr := shared.Validate(body); valErr != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(web.BadRequest(valErr))
	}

	db := database.DbCtx()

	newEntity := database.Product{
		Name:  body.Name,
		Price: body.Price,
	}

	addErr := db.Products().Add(&newEntity)
	if addErr != nil {
		log.Fatal(addErr)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(addErr)
	}

	return ctx.JSON(Product{
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
func get(ctx *fiber.Ctx) error {

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	db := database.DbCtx()

	product := db.Products().GetById(id)
	if product == nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).
			JSON(web.NotFound())
	}

	return ctx.JSON(Product{
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
func remove(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	db := database.DbCtx()

	err = db.Products().Remove(id)
	if err != nil {
		log.Fatal(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func ConfigureRoutes(app *fiber.App) {
	products := app.Group("/products", web.Auth())
	products.Get("/", getAll)
	products.Post("/", add)
	products.Get("/:id", get)
	products.Delete("/:id", remove)
}
